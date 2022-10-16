package infrastructure

import (
	"service-auth/pkg/model"
)

func FetchUserByLogin(login string) *model.User {
	user := new(model.User)
	user.Id = "123"
	user.Token = "old"
	user.Login = login
	user.Password = "doe"
	return user
}

func RegenerageToken(user *model.User) {
	(*user).Token = "foobar"
}

func FetchUserByToken(token string) *model.User {
	user := new(model.User)
	user.Id = "123"
	user.Token = token
	user.Login = "login"
	user.Password = "doe"
	return user
}
