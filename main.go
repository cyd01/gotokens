package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gotokens/flags"
	"gotokens/tokens"

	"github.com/gin-gonic/gin"
)

var (
	f        = flags.NewFlag("tokens")
	addr     = f.String("addr", ":80", "bind address")
	dir      = f.String("dir", ".", "root directory")
	expire   = f.Int("expire", 300, "expiration time (seconds)")
	login    = f.String("login", "admin", "admin login")
	password = f.String("password", "pass", "admin password")
)

// Main procedure
func main() {
	// Reading command-line flags
	f.Parse(os.Args[f.NArg()+1:])
	tokens.AddTokenUser(*login, *password)

	tokens.TokensSetExpirationTime(*expire)

	rand.Seed(time.Now().UnixNano())
	/* Switch to production mode
	   - using env:   export GIN_MODE=release
	   - using code:  gin.SetMode(gin.ReleaseMode)
	*/

	// Setting routes for api
	router := gin.Default()

	// Serve alive service
	router.GET("/alive", func(c *gin.Context) { c.JSON(200, gin.H{"status": "success", "message": "alive"}) })

	TokensGroup := router.Group("/tokens")
	{
		TokensGroup.GET("/challengedata", tokens.TokensGetChallengeData)
		TokensGroup.GET("/", tokens.TokensGet)             /* with auth */
		TokensGroup.POST("/clean", tokens.TokensPostClean) /* with auth */
		TokensGroup.GET("/validate/:token", tokens.TokensGetValidate)
		TokensGroup.GET("/:id", tokens.TokensGetId)       /* with auth */
		TokensGroup.DELETE("/:id", tokens.TokensDeleteId) /* with auth */
		TokensGroup.POST("/", tokens.TokensPost)
		TokensGroup.POST("/auth", tokens.TokensPostAuth)
		TokensGroup.GET("/admin.html", func(c *gin.Context) { c.File(*dir + "/admin.html") })
	}

	router.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "admin.html" {
			c.File(*dir + "/admin.html")
		} else if id == "favicon.ico" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Not found"})
		}
	})
	router.GET("/", func(c *gin.Context) { c.File(*dir + "/admin.html") })

	// Define the server
	if !strings.HasPrefix(*addr, ":") {
		*addr = ":" + *addr
	}
	srv := &http.Server{
		Addr:    *addr,
		Handler: router,
	}

	log.Println("Starting server " + *addr + " with expiration date at " + strconv.Itoa(*expire) + " seconds")

	// Starting
	go func() {
		if err := router.Run(*addr); err != nil {
			log.Fatalf("Server can't start: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
