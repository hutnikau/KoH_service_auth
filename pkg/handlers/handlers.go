package handlers

import (
	"net/http"

	"service-auth/pkg/infrastructure"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorBody struct {
	ErrorMsg string `json:"error,omitempty"`
}

var ErrorMethodNotAllowed = "method Not allowed"

func Authenticate(login string, password string) (*events.APIGatewayProxyResponse, error) {
	user := infrastructure.FetchUserByLogin(login)

	if user.IsPasswordValid(password) {
		infrastructure.RegenerageToken(user)
		return apiResponse(http.StatusOK, user)
	}

	return apiResponse(http.StatusBadRequest, ErrorBody{"User not found"})
}

func VerifyToken(token string) (*events.APIGatewayProxyResponse, error) {
	user := infrastructure.FetchUserByToken(token)

	if user != nil {
		return apiResponse(http.StatusOK, user)
	}

	return apiResponse(http.StatusBadRequest, ErrorBody{"User not authorized"})
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
