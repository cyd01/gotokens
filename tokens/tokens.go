package tokens

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gotokens/tools"

	"github.com/gin-gonic/gin"
)

/* The time-to-live for token */
const (
	tokenTTL = 300
)

/* Expiration time (in seconds) */
var (
	expireTime int = 300
	TokenCode      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func TokensSetExpirationTime(ex int) {
	expireTime = ex
}

/* The challenge data properties */
type CHALLENGEDATA struct {
	Id      string `json:"-"`
	Data    string `json:"challengedata"`
	Created int64  `json:"-"`
}

/* The challenge data database */
var ChallengeData []CHALLENGEDATA

/* Get a challenge data (GET /challengedata)
 */
func TokensGetChallengeData(c *gin.Context) {
	data := ""
	now := tools.Epoch()
	for i := 0; i < 16; i++ {
		data = data + tools.Gensha256(strconv.FormatInt(now+int64(i), 10))
	}
	id := tools.Genuuid()
	item := CHALLENGEDATA{
		Id:      id,
		Data:    data,
		Created: now,
	}
	ChallengeData = append(ChallengeData, item)
	c.SetCookie("ChallengeData", id, tokenTTL, "/", tools.Replace(":[0-9]*$", "", c.Request.Host), false, true)
	c.JSON(http.StatusOK, item)
}

/* The token properties */
type TOKEN struct {
	Id      string `json:"id"`
	User    string `json:"user"`
	Token   string `json:"token"`
	Address string `json:"address"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
	Hits    int64  `json:"hits"`
}

/* The tokens database */
var Tokens []TOKEN

/* The users list that are authorized to create a token : map[login] => password */
var tokenUsers map[string]string

func AddTokenUser(login, password string) {
	tokenUsers[login] = password
}

/* We read the users list from file */
func init() {
	TokenCode = tools.Shuffle(TokenCode)
	tokenUsers = make(map[string]string)
	tools.ReadFromJSONFile("users.json", &tokenUsers)
}

/* Clean token and challenge data database on expiration date */
func TokensClean() {
	now := tools.Epoch()
	if len(Tokens) > 0 {
		for i := 0; i < len(Tokens); i++ {
			if (Tokens[i].Updated + int64(expireTime)) < now {
				log.Println("Remove token " + Tokens[i].Token)
				Tokens = append(Tokens[:i], Tokens[i+1:]...)
			}
		}
	}
	if len(ChallengeData) > 0 {
		for i := 0; i < len(ChallengeData); i++ {
			if (ChallengeData[i].Created + int64(expireTime)) < now {
				log.Println("Remove challentge data " + ChallengeData[i].Data)
				ChallengeData = append(ChallengeData[:i], ChallengeData[i+1:]...)
			}
		}
	}
}

/* Validate a given userToken (see TestToken func below)
 * A userToken is in the form user-token
 * return is true => the token is valid
 * return is false => the token is invalid or unknown
 */
func TokensValidate(userToken string) bool {
	test := false
	now := tools.Epoch()
	userTokenSplit := strings.Split(userToken, "-")
	if len(userTokenSplit) != 2 {
		return false
	}
	user, _ := tools.StringDecode(userTokenSplit[0], TokenCode)
	token := userTokenSplit[1]
	if len(Tokens) > 0 {
		for i := 0; i < len(Tokens); i++ {
			if (user == Tokens[i].User) && (token == Tokens[i].Token) {
				if (Tokens[i].Updated + int64(expireTime)) < now {
					log.Println("Remove token " + Tokens[i].Token)
					Tokens = append(Tokens[:i], Tokens[i+1:]...)
				} else {
					test = true
					Tokens[i].Hits = Tokens[i].Hits + 1
					Tokens[i].Updated = now
					log.Println("Token validated for user " + user)
					break
				}
			}
		}
	}
	if !test {
		log.Println("Token is not valid")
	}
	return test
}

/* Test the userToken received from client (in query, cookie or header)
 * A userToken is in the form user-token
 * The userToken is passed to TokensValidate func above
 */
func TestToken(c *gin.Context) bool {
	test := false
	userToken, b := c.GetQuery("token")
	if !b {
		userToken = ""
	}
	if len(userToken) == 0 {
		if s, err := c.Cookie("Token"); err == nil {
			userToken = s
		}
	}
	if len(userToken) == 0 {
		userToken = c.GetHeader("TOKEN")
	}
	if len(userToken) > 0 {
		test = TokensValidate(userToken)
	}
	return test
}

/* API */

/* Get all the tokens (GET /tokens)
 * with auth
 * 401 -> Unauthorized
 * 200 -> Ok
 */
func TokensGet(c *gin.Context) {
	if !TestToken(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, Tokens)
}

/* Get one token (GET /tokens/:id)
 * with auth
 * 401 -> Unauthorized
 * 404 -> Not found
 * 200 -> Ok
 */
func TokensGetId(c *gin.Context) {
	if !TestToken(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Unauthorized"})
		return
	}
	id := c.Param("id")
	now := tools.Epoch()
	if len(Tokens) > 0 {
		for i := 0; i < len(Tokens); i++ {
			if id == Tokens[i].Id {
				if (Tokens[i].Updated + int64(expireTime)) < now {
					log.Println("Remove token " + Tokens[i].Token)
					Tokens = append(Tokens[:i], Tokens[i+1:]...)
				} else {
					Tokens[i].Hits = Tokens[i].Hits + 1
					Tokens[i].Updated = now
					c.JSON(http.StatusOK, Tokens[i])
					return
				}
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Not found"})
}

/* Delete one token (DELETE /tokens/:id)
 * with auth
 * 401 -> Unauthorized
 * 404 -> Not found
 * 204 -> Deleted
 */
func TokensDeleteId(c *gin.Context) {
	if !TestToken(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Unauthorized"})
		return
	}
	id := c.Param("id")
	if len(Tokens) > 0 {
		for i := 0; i < len(Tokens); i++ {
			if id == Tokens[i].Id {
				log.Println("Remove token " + Tokens[i].Token)
				Tokens = append(Tokens[:i], Tokens[i+1:]...)
				c.Status(http.StatusNoContent)
				return
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Not found"})
}

func TokensSetCookie(c *gin.Context, login, token string) {
	c.SetCookie("Token", tools.StringEncode(login, token)+"-"+token, tokenTTL, "/", tools.Replace(":[0-9]*$", "", c.Request.Host), false, true)
}

/* Validate one token (GET /validate/:token)
 * no auth
 * 200 -> Ok
 * 404 -> Not found or invalid
 */
func TokensGetValidate(c *gin.Context) {
	token := c.Param("token")
	test := false
	if len(Tokens) > 0 {
		for i := 0; i < len(Tokens); i++ {
			if token == Tokens[i].Token {
				test = true
				break
			}
		}
	}
	if test {
		TokensSetCookie(c, "Unknown", token)
		//c.SetCookie("Token", tools.StringEncode("Unknown", token)+"-"+token, tokenTTL, "/", tools.Replace(":[0-9]*$", "", c.Request.Host), false, true)
		log.Println("Token is valid")
		c.JSON(http.StatusOK, gin.H{"status": "succeeded", "message": "Valid token"})
	} else {
		log.Println("Token is not valid")
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Not found"})
	}
}

/* The credential structure for user/password */
type INPUTCREDENTIALS struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/* Create a new token (POST /tokens) for a user with credentials in request body {"login":"xxx","password":"yyy"}
 * no auth
 * 400 -> Wrong parameter
 * 401 -> Wrong credentials
 * 201 -> Token created (cookie post)
 */
func TokensPost(c *gin.Context) {
	var input INPUTCREDENTIALS
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	var challengeData string = ""
	if id, err := c.Cookie("ChallengeData"); err == nil {
		if len(ChallengeData) > 0 {
			for i := 0; i < len(ChallengeData); i++ {
				if ChallengeData[i].Id == id {
					challengeData = ChallengeData[i].Data
				}
			}
		}
	}
	if len(challengeData) > 0 {
		data := fmt.Sprintf("%x", md5.Sum([]byte(tokenUsers[input.Login]+challengeData)))
		if data != input.Password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Unauthorized"})
			return
		}
	} else if tokenUsers[input.Login] != input.Password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Unauthorized"})
		return
	}
	item := GenerateToken(input.Login, c.Request.RemoteAddr)
	c.SetCookie("ChallengeData", "", -1, "/", tools.Replace(":[0-9]*$", "", c.Request.Host), false, true)
	c.SetCookie("Token", tools.StringEncode(item.User, TokenCode)+"-"+item.Token, tokenTTL, "/", tools.Replace(":[0-9]*$", "", c.Request.Host), false, true)
	c.JSON(http.StatusCreated, item)
}

/* Create a new token (POST /tokens/auth) for a user with credentials basic auth
 * no auth
 * 204 -> already connected
 * 401 -> Wrong credentials
 * 201 -> Token created (cookie post)
 */
func TokensPostAuth(c *gin.Context) {
	if !TestToken(c) {
		user, pass, hasAuth := c.Request.BasicAuth()
		if hasAuth && tokenUsers[user] == pass {
			item := GenerateToken(user, c.Request.RemoteAddr)
			c.SetCookie("Token", tools.StringEncode(item.User, TokenCode)+"-"+item.Token, tokenTTL, "/", tools.Replace(":[0-9]*$", "", c.Request.Host), false, true)
			c.JSON(http.StatusCreated, item)
		} else {
			c.Status(http.StatusUnauthorized)
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		}
	} else {
		c.Status(http.StatusNoContent)
	}
}

/* Clean token (POST /tokens/clean)
 * with auth
 * 401 -> Unauthorized
 * 204 -> Cleaned
 */
func TokensPostClean(c *gin.Context) {
	if !TestToken(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Unauthorized"})
		return
	}
	TokensClean()
	c.Status(http.StatusNoContent)
}

/* Function to generate a new token */
func GenerateToken(user string, RemoteAddr string) TOKEN {
	id := tools.Genuuid()
	address := tools.Replace(":[^:]*$", "", RemoteAddr)
	now := tools.Epoch()
	token := tools.Gensha256(id + "/" + strconv.FormatInt(now, 10))
	item := TOKEN{
		Id:      id,
		User:    user,
		Token:   token,
		Address: address,
		Created: now,
		Updated: now,
		Hits:    0,
	}
	log.Println("Create token " + token + " for user " + user)
	Tokens = append(Tokens, item)
	return item
}
