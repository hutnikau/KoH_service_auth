package handlers

import "service-auth/pkg/model"

type TokenRepository interface {
	RegenerateToken(user *model.User) *model.Token
	FetchUserIdByToken(token string) (string, error)
}

type UserRepository interface {
	FetchUserById(userId string) (model.User, error)
	FetchUserByLogin(login string) (model.User, error)
}
