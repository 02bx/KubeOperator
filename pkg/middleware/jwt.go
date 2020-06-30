package middleware

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/spf13/viper"
	"time"
)

var (
	secretKey []byte
	exp       int
)

func JWTMiddleware() *jwtmiddleware.Middleware {
	secretKey = []byte(viper.GetString("jwt.secret"))
	exp = viper.GetInt("jwt.exp")
	return jwtmiddleware.New(jwtmiddleware.Config{
		Extractor: jwtmiddleware.FromAuthHeader,
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//自己加密的秘钥或者说盐值
			return secretKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler:  ErrorHandler,
	})
}

func ErrorHandler(ctx context.Context, err error) {
	if err == nil {
		return
	}
	ctx.StopExecution()
	response := &dto.Response{
		Msg: err.Error(),
	}
	ctx.StatusCode(iris.StatusInternalServerError)
	ctx.JSON(response)
}

func LoginHandler(ctx context.Context) {
	aul := new(auth.Credential)
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(dto.Response{Msg: err.Error()})
		return
	}

	data, err := CheckLogin(aul.Username, aul.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(dto.Response{Msg: err.Error()})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(data)
	return
}

func CheckLogin(username string, password string) (*auth.JwtResponse, error) {
	user, err := service.UserAuth(username, password)
	if err != nil {
		return nil, err
	}
	token, err := CreateToken(user)
	if err != nil {
		return nil, err
	}
	resp := new(auth.JwtResponse)
	resp.Token = token
	resp.User = *user
	return resp, err
}

func CreateToken(user *auth.SessionUser) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     user.Name,
		"email":    user.Email,
		"userId":   user.UserId,
		"isActive": user.IsActive,
		"language": user.Language,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute * time.Duration(exp)).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func GetAuthUser(ctx context.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)
	foobar := user.Claims.(jwt.MapClaims)
	sessionUserJson, _ := json.Marshal(foobar)
	sessionUserJsonStr := string(sessionUserJson)
	var sessionUser auth.SessionUser
	json.Unmarshal([]byte(sessionUserJsonStr), &sessionUser)
	resp := new(auth.JwtResponse)
	resp.User = sessionUser
	resp.Token = user.Raw
	ctx.JSON(resp)
}
