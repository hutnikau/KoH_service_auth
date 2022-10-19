package handlers

import (
	"net/http"
	"service-auth/pkg/infrastructure"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorBody struct {
	ErrorMsg string `json:"error,omitempty"`
}

func Authenticate(login string, password string) (*events.APIGatewayV2HTTPResponse, error) {
	userRepo := infrastructure.NewUserRepository()
	tokenRepo := infrastructure.NewTokenRepository()
	user, err := userRepo.FetchUserByLogin(login)

	if err != nil {
		return apiResponse(http.StatusNotFound, ErrorBody{"User not found"})
	}

	if user.IsPasswordValid(password) {
		tokenRepo.RegenerageToken(&user)
		user.Password = ""
		return apiResponse(http.StatusOK, user)
	}

	return apiResponse(http.StatusForbidden, ErrorBody{"Wrong login or password"})
}

func VerifyToken(token string) (*events.APIGatewayV2HTTPResponse, error) {
	tokenRepo := infrastructure.NewTokenRepository()
	userId, err := tokenRepo.FetchUserIdByToken(token)

	if err == nil {
		return apiResponse(http.StatusOK, userId)
	}

	return apiResponse(http.StatusNotFound, ErrorBody{"User not authorized"})
}
