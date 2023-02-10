package auth

import (
    "UserService/utils"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
    "time"
)

type (
    Claims struct {
        UID  string
        Data map[string]interface{}
        jwt.StandardClaims
    }
    AuthorizedRequestFunc func(c *gin.Context)
    GenerateTokenFunc     func(UID string, duration time.Duration, data map[string]interface{}) (time.Time, string, error)
)

func New(JwtKey string) (gin.HandlerFunc, GenerateTokenFunc) {
    return func(c *gin.Context) {
            claims, ok := isAuthorizedRequest(JwtKey, c.Request)
            if !ok {
                c.AbortWithStatusJSON(http.StatusOK, gin.H{
                    "code": -1,
                    "msg":  "unauthorized",
                })
            } else {
                c.Set("uid", claims.UID)
                c.Set("data", claims.Data)
            }
        }, func(UID string, duration time.Duration, data map[string]interface{}) (time.Time, string, error) {
            return generateJWTToken(JwtKey, UID, duration, data)
        }
}
func isAuthorizedRequest(JwtKey string, r *http.Request) (*Claims, bool) {
    reqToken := r.Header.Get("Authorization")
    splitToken := strings.Split(reqToken, "Bearer")
    if len(splitToken) != 2 {
        return nil, false
    }
    tokenString := strings.TrimSpace(splitToken[1])
    claims := &Claims{
        Data: map[string]interface{}{},
    }
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(JwtKey), nil
    })
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return nil, false
        }
        return nil, false
    }
    if !token.Valid {
        return nil, false
    }
    return claims, true
}
func generateJWTToken(JwtKey string, UID string, duration time.Duration, data map[string]interface{}) (time.Time, string, error) {
    expirationTime := utils.GetTime().Add(duration)
    claims := &Claims{
        UID:  UID,
        Data: data,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(JwtKey))
    return expirationTime, tokenString, err
}
