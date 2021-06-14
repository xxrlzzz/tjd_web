package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/gredis"
	"traffic_jam_direction/pkg/util"

	"github.com/gin-gonic/gin"

	"traffic_jam_direction/pkg/e"
)

type TokenType struct {
	Token string `json:"token"`
}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	// JWT(JSON WEB TOKEN)
	// 检查 token 和是否过期
	return func(c *gin.Context) {
		code := e.SUCCESS

		appG := app.Gin{C: c}
		var token TokenType

		for {
			if err := c.ShouldBindWith(&token, binding.JSON); err != nil || token.Token == "" {
				code = e.ErrorAuthTokenNotExist
				break
			}
			// 用户已经退出登录
			if gredis.Exists("token" + token.Token) {
				code = e.ErrorAuthTokenLogout
				break
			}
			cli, err := util.ParseToken(token.Token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ErrorAuthCheckTokenTimeout
				default:
					code = e.ErrorAuthCheckTokenFail
				}
				break
			}

			// 设置用户登录状态
			c.Set("login", cli.Login)
		}

		if code != e.SUCCESS {
			appG.Response(http.StatusUnauthorized, code, nil)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
