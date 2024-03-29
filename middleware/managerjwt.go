package jwt

import (
	"net/http"
	"simple_project/pkg/setting"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"

	"simple_project/pkg/e"
	"simple_project/pkg/utils"
)

// JWT is jwt middleware
func ManagerJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader(setting.AppSetting.AuthKey)
		if token == "" {
			code = e.NO_AUTH
		} else {
			_, err := utils.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
