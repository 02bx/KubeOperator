package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/user"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
)

var (
	UserNotFound     = errors.New("can not find user")
	PasswordNotMatch = errors.New("password not match")
)

func UserAuth(name string, password string) (sessionUser *auth.SessionUser, err error) {
	var dbUser user.User
	if db.DB.Where("name = ?", name).First(&dbUser).RecordNotFound() {
		return nil, UserNotFound
	}
	password, err = encrypt.StringEncrypt(password)
	if err != nil {
		return nil, err
	}
	if dbUser.Password != password {
		return nil, PasswordNotMatch
	}
	return dbUser.ToSessionUser(), nil
}
