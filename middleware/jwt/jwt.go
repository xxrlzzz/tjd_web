package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"traffic_jam_direction/pkg/util"

	"github.com/gin-gonic/gin"

	"traffic_jam_direction/pkg/e"
)

type Token struct {
	Token string `json:"token"`
}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	// JWT(JSON WEB TOKEN)
	// 检查 token 和是否过期
	return func(c *gin.Context) {
		code := e.SUCCESS
		//var data interface{}

		//token := c.Query("token")
		var Token Token
		if err := c.ShouldBindWith(&Token, binding.JSON); err != nil {
			code = e.ERROR_AUTH_TOKEN_NOT_EXIST
		} else {
			token := Token.Token
			if token == "" {
				code = e.ERROR_AUTH_TOKEN_NOT_EXIST
			} else {
				cli, err := util.ParseToken(token)
				if err != nil {
					switch err.(*jwt.ValidationError).Errors {
					case jwt.ValidationErrorExpired:
						code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
					default:
						code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
					}
				} else {
					// 设置用户登录状态
					c.Set("login", cli.Login)
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": "",
			})
			c.Abort()
		} else {
			c.Next()
		}
	}
}
