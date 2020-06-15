package middleware

import (
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
)

var log = logger.Default
var sessionUser = auth.SessionUser{}

func JWTMiddleware() *jwtmiddleware.Middleware {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//自己加密的秘钥或者说盐值
			return []byte("My Secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func ErrorHandler() {

}

func LoginHandler(ctx iris.Context) {
	aul := new(auth.Credential)
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(dto.Response{Status: false, Msg: "", Data: "请求参数错误"})
		return
	}

	data, err := CheckLogin(aul.Username, aul.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(dto.Response{Status: false, Msg: "校验失败", Data: nil})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(dto.Response{Msg: "success", Data: data, Status: true})
	return
}

func CheckLogin(username string, password string) (*auth.JwtResponse, error) {
	user, err := service.UserAuth(username, password)
	if err != nil {
		return nil, err
	}
	token, err := CreateJWTToken(user)
	if err != nil {
		return nil, err
	}
	resp := new(auth.JwtResponse)
	resp.Token = token
	resp.User = *user
	return resp, err
}

func CreateJWTToken(user *auth.SessionUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"name":     user.Name,
		"email":    user.Email,
		"language": user.Language,
		"isActive": user.IsActive,
		"userId":   user.UserId,
	})
	tokenString, err := token.SignedString([]byte("My Secret"))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func myAuthenticatedHandler(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)

	ctx.Writef("This is an authenticated request\n")
	ctx.Writef("Claim content:\n")

	foobar := user.Claims.(jwt.MapClaims)
	for key, value := range foobar {
		ctx.Writef("%s = %s", key, value)
	}
}

//func JWTMiddleware() *jwt.GinJWTMiddleware {
//	secret := viper.GetString("app.secret")
//	j, err := jwt.New(&jwt.GinJWTMiddleware{
//		Key:           []byte(secret),
//		Timeout:       time.Hour,
//		MaxRefresh:    time.Hour,
//		TimeFunc:      time.Now,
//		TokenHeadName: "Bearer",
//		IdentityKey:   "user",
//		Authenticator: func(ctx *gin.Context) (i interface{}, err error) {
//			var credential auth.Credential
//			if err := ctx.ShouldBind(&credential); err != nil {
//				return nil, jwt.ErrMissingLoginValues
//			}
//			sUser, err := service.UserAuth(credential.Username, credential.Password)
//			if err != nil {
//				if sUser != nil && sUser.IsActive == false {
//					return nil, err
//				} else {
//					return nil, err
//				}
//			}
//			return sUser, nil
//		},
//		PayloadFunc: func(data interface{}) jwt.MapClaims {
//			if v, ok := data.(*auth.SessionUser); ok {
//				sessionUser = *v
//				return jwt.MapClaims{
//					"user": v,
//				}
//			}
//			return jwt.MapClaims{}
//		},
//		IdentityHandler: func(ctx *gin.Context) interface{} {
//			claims := jwt.ExtractClaims(ctx)
//			return &claims
//		},
//		LoginResponse: func(ctx *gin.Context, code int, token string, expire time.Time) {
//			ctx.JSON(http.StatusOK, gin.H{
//				"code":   code,
//				"token":  token,
//				"expire": expire.Format(time.RFC3339),
//				"user":   sessionUser,
//			})
//			return
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	return j
//}
//
//func GetAuthUser(ctx *gin.Context) {
//	claims := jwt.ExtractClaims(ctx)
//	token := ctx.Keys["JWT_TOKEN"].(string)
//	ctx.JSON(http.StatusOK, gin.H{
//		"code":  http.StatusOK,
//		"token": token,
//		"user":  claims["user"],
//	})
//	return
//}
