package middleware

import (
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

var log = logger.Default

func JWTMiddleware() *jwt.GinJWTMiddleware {
	secret := viper.GetString("app.secret")
	j, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:           []byte(secret),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TimeFunc:      time.Now,
		TokenHeadName: "Bearer",
		IdentityKey:   "user",
		Authenticator: func(ctx *gin.Context) (i interface{}, err error) {
			var credential auth.Credential
			if err := ctx.ShouldBind(&credential); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			user, err := service.UserAuth(credential.Username, credential.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return user, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*auth.SessionUser); ok {
				return jwt.MapClaims{
					"name":   v.Name,
					"userId": v.UserId,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx *gin.Context) interface{} {
			claims := jwt.ExtractClaims(ctx)
			return &auth.SessionUser{
				UserId: claims["userId"].(string),
				Name:   claims["name"].(string),
			}
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return j
}
