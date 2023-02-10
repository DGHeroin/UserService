package UserAPI

import (
    "UserService/config"
    "UserService/dao"
    "UserService/logger"
    "UserService/models"
    "UserService/utils"
    "UserService/validator"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/patrickmn/go-cache"
    "math/rand"
    "net/http"
    "time"
)

var (
    tc *cache.Cache
)

func init() {
    tc = cache.New(time.Minute*10, time.Minute)
}

func HandleAllocEmailVerifyCode(sendMailFunc func(string) error) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request = &struct {
            Email string `json:"email"`
        }{}
        if err := c.BindJSON(request); err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 10001,
                "msg":  "bad request",
            })
            return
        }
        logger.Println(request)
        if !validator.VerifyEmailFormat(request.Email) {
            c.JSON(http.StatusOK, gin.H{
                "code": 10002,
                "msg":  "param invalid",
            })
            return
        }

        verifyCode := fmt.Sprint(rand.Intn(899999) + 100000)

        err := sendMailFunc(verifyCode)

        if !validator.VerifyEmailFormat(request.Email) {
            c.JSON(http.StatusOK, gin.H{
                "code": 10002,
                "msg":  "param invalid",
            })
            logger.Println("send mail error:", err)
            return
        }

        tc.SetDefault(request.Email, map[string]interface{}{
            "time": utils.GetUnixTime(),
            "code": verifyCode,
        })

        logger.Println("申请验证码", request.Email, verifyCode)

        c.JSON(http.StatusOK, gin.H{
            "code": 0,
        })
    }
}
func HandleUserRegister(tokenFunc func(c *gin.Context, user *models.User) (expire time.Time, token string, err error)) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request = &struct {
            Email    string `json:"email"`
            Password string `json:"password"`
            Code     string `json:"code"`
        }{}
        if err := c.BindJSON(request); err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 10001,
                "msg":  "bad request",
            })
            return
        }
        if config.UseEmailVerify { // 验证 email
            if p, ok := tc.Get(request.Email); !ok {
                c.JSON(http.StatusOK, gin.H{
                    "code": 20002,
                    "msg":  "verify code expired",
                })
                return
            } else {
                mm := p.(map[string]interface{})
                if mm["code"].(string) != request.Code {
                    c.JSON(http.StatusOK, gin.H{
                        "code": 20003,
                        "msg":  "verify code invalid",
                    })
                    return
                }
            }

            // 验证通过
            tc.Delete(request.Email)
        }
        // 执行注册
        user, err := models.AddUser(dao.SharedDB(), request.Email, request.Password)
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 50003,
                "msg":  "internal error",
            })
            logger.Println("[User] register error:", err)
            return
        }
        expire, tokenString, err := tokenFunc(c, user)
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 50006,
                "msg":  "internal error",
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "code":   0,
            "token":  tokenString,
            "expire": expire.Unix(),
        })
    }
}
func HandleUserLogin(tokenFunc func(c *gin.Context, user *models.User) (expire time.Time, token string, err error)) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request = &struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }{}
        if err := c.BindJSON(request); err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 10001,
                "msg":  "bad request",
            })
            return
        }

        // 执行登录
        user, err := models.GetUsersByEmail(dao.SharedDB(), request.Email)
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 50003,
                "msg":  "internal error",
            })
            logger.Println("[User] get user by email error:", err)
            return
        }
        if ok, err := utils.ComparePasswordAndHash(request.Password, user.Password); err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 50004,
                "msg":  "internal error",
            })
            logger.Println(err)
            return
        } else {
            if !ok {
                c.JSON(http.StatusOK, gin.H{
                    "code": 50005,
                    "msg":  "login fail",
                })
                return
            }
        }
        expire, tokenString, err := tokenFunc(c, user)
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 50006,
                "msg":  "internal error",
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "code":   0,
            "token":  tokenString,
            "expire": expire.Unix(),
        })
    }
}
func HandleUserRefreshToken(tokenFunc func(c *gin.Context, user *models.User) (expire time.Time, token string, err error)) gin.HandlerFunc {
    return func(c *gin.Context) {
        uid := c.Value("uid")
        user, err := models.GetUsersByUID(dao.SharedDB(), uid.(string))
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 50003,
                "msg":  "internal error",
            })
            logger.Println("[User] get user by email error:", err)
            return
        }
        expire, tokenString, err := tokenFunc(c, user)
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 50006,
                "msg":  "internal error",
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "code":   0,
            "token":  tokenString,
            "expire": expire.Unix(),
        })
    }
}
