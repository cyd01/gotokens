package gotokens

import (
    "log"
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
)

/* The time-to-live for token */
const (
    tokenTTL = 300
)

/* The token properties */
type TOKENS struct {
    Id           string            `json:"id"`
    User         string            `json:"user"`
    Token        string            `json:"token"`
    Address      string            `json:"address"`
    Created      int64             `json:"created"`
    Updated      int64             `json:"updated"`
    Hits         int64             `json:"hits"`
}

/* The tokens database */
var Tokens []TOKENS

/* The users list that are authorized to create a token : map[login] => password */
var tokenUsers map[string]string

/* We read the users list from file */
func init() {
    tokenUsers = make(map[string]string)
    ReadFromJSON("users.json",&tokenUsers)
}

/* Clean token database on expiration date */
func TokensClean() {
    now := epoch()
    if( len(Tokens)>0 ) {
        for i:=0; i<len(Tokens); i++ {
            if( (Tokens[i].Updated+int64(*expire)) < now ) {
                log.Println("Remove token "+Tokens[i].Token)
                Tokens = append(Tokens[:i], Tokens[i+1:]...)
            }
        }
    }
}

/* Validate a given userToken (see TestToken func below)
 * A userToken is in the form user-token
 * return is true => the token is valid
 * return is false => the token is invalid or unknown
 */
func TokensValidate( userToken string ) bool {
    test := false
    now := epoch()
    userTokenSplit := strings.Split(userToken,"-")
    if( len(userTokenSplit)!=2 ) { return false }
    user := userTokenSplit[0]
    token := userTokenSplit[1]
    if( len(Tokens)>0 ) {
        for i:=0; i<len(Tokens); i++ {
            if( (user==Tokens[i].User) && (token==Tokens[i].Token) ) {
                if( (Tokens[i].Updated+int64(*expire)) < now ) {
                    log.Println("Remove token "+Tokens[i].Token)
                    Tokens = append(Tokens[:i], Tokens[i+1:]...)
                } else {
                    test = true
                    Tokens[i].Hits = Tokens[i].Hits + 1
                    Tokens[i].Updated = now
                    break
                }
            }
        }
    }
    return test
}

/* Test the userToken received from client (in query, cookie or header)
 * A userToken is in the form user-token
 * The userToken is passed to TokensValidate func above
 */
func TestToken( c *gin.Context ) bool {
    test := false
    userToken, b := c.GetQuery("token")
    if( !b ) { userToken = "" }
    if( len(userToken)==0 ) {
        if s, err := c.Cookie("Token"); err == nil { userToken = s }
    }
    if( len(userToken)==0 ) {
        userToken = c.GetHeader("TOKEN")
    }
    if( len(userToken)>0 ) {
        test = TokensValidate( userToken )
    }
    return test
}

/* API */

/* Get all the tokens (GET /tokens)
 * with auth
 * 401 -> Unauthorized
 * 200 -> Ok
 */
func GetTokens(c *gin.Context) { 
    if( !TestToken(c) ) {
        c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"status": "failed", "message": "Unauthorized"})
        return
    }
    c.JSON(http.StatusOK,Tokens)
}

/* Get one token (GET /tokens/:id)
 * with auth
 * 401 -> Unauthorized
 * 404 -> Not found
 * 200 -> Ok
 */
func GetTokensId(c *gin.Context) {
    if( !TestToken(c) ) {
        c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"status": "failed", "message": "Unauthorized"})
        return
    }
    id := c.Param("id")
    now := epoch()
    if( len(Tokens)>0 ) {
        for i:=0; i<len(Tokens); i++ {
            if( id == Tokens[i].Id ) {
                if( (Tokens[i].Updated+int64(*expire)) < now ) {
                    log.Println("Remove token "+Tokens[i].Token)
                    Tokens = append(Tokens[:i], Tokens[i+1:]...)
                } else {
                    Tokens[i].Hits = Tokens[i].Hits + 1
                    Tokens[i].Updated = now
                    c.JSON(http.StatusOK,Tokens[i])
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
func DeleteTokensId(c *gin.Context) {
    if( !TestToken(c) ) {
        c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"status": "failed", "message": "Unauthorized"})
        return
    }
    id := c.Param("id")
    if( len(Tokens)>0 ) {
        for i:=0; i<len(Tokens); i++ {
            if( id == Tokens[i].Id ) {
                log.Println("Remove token "+Tokens[i].Token)
                Tokens = append(Tokens[:i], Tokens[i+1:]...)
                c.Status( http.StatusNoContent )
                return
            }
        }
    }
    c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Not found"})
}

/* Validate one token (GET /validate/:token)
 * no auth
 * 200 -> Ok
 * 404 -> Not found or invalid
 */
func GetValidateToken(c *gin.Context) {
    userToken := c.Param("token")
    test := TokensValidate(userToken)
    if( test ) {
        c.JSON( http.StatusOK, gin.H{"status": "succeeded", "message": "Valid token"} )
    } else {
        c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Not found"})
    }
}

/* The credential structure for user/password */
type INPUTCREDENTIALS struct {
    Login        string            `json:"login" binding:"required"`
    Password     string            `json:"password" binding:"required"`
}

/* Create a new token (POST /tokens) for a user with credentials in request body {"login":"xxx","password":"yyy"}
 * no auth
 * 400 -> Wrong parameter
 * 401 -> Wrong credentials
 * 201 -> Token created (cookie post)
 */
func PostTokens(c *gin.Context) {
    var input INPUTCREDENTIALS
    if err := c.BindJSON(&input); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
        return
    }
    if tokenUsers[input.Login] != input.Password {
        c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"status": "failed", "message": "Unauthorized"})
        return
    }
    item := GenerateToken(input.Login,c.Request.RemoteAddr)
    c.SetCookie("Token",item.User+"-"+item.Token,tokenTTL,"/",replace(":[0-9]*$","",c.Request.Host),false,true)
    c.JSON(http.StatusCreated,item)
}

/* Create a new token (POST /tokens/auth) for a user with credentials basic auth
 * no auth
 * 204 -> already connected
 * 401 -> Wrong credentials
 * 201 -> Token created (cookie post)
 */
 func PostAuthTokens(c *gin.Context) {
    if( !TestToken(c) ) {
        user, pass, hasAuth := c.Request.BasicAuth()
        if hasAuth && tokenUsers[user] == pass {
            item := GenerateToken(user,c.Request.RemoteAddr)
            c.SetCookie("Token",item.User+"-"+item.Token,tokenTTL,"/",replace(":[0-9]*$","",c.Request.Host),false,true)
            c.JSON(http.StatusCreated,item)
        } else {
            c.Status( http.StatusUnauthorized )
            c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
        }
    } else {
        c.Status( http.StatusNoContent )
    }
}

/* Clean token (POST /tokens/clean) 
 * with auth
 * 401 -> Unauthorized
 * 204 -> Cleaned
 */
func PostCleanTokens(c *gin.Context) {
    if( !TestToken(c) ) {
        c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"status": "failed", "message": "Unauthorized"})
        return
    }
    TokensClean()
    c.Status( http.StatusNoContent )
}

/* Function to generate a new token */
func GenerateToken(user string, RemoteAddr string) TOKENS {
    id := genuuid()
    address := replace(":[^:]*$","",RemoteAddr)
    now := epoch()
    token := gensha256(id+"/"+strconv.FormatInt(now,10))
    item := TOKENS{
        Id:         id,
        User:       user,
        Token:      token,
        Address:    address,
        Created:    now,
        Updated:    now,
        Hits:       0,
    }
    log.Println("Create token "+token+" for user "+user)
    Tokens = append(Tokens, item)
    return item
}
