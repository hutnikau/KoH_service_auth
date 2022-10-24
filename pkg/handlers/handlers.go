package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorBody struct {
	ErrorMsg string `json:"error,omitempty"`
}

type Handler struct {
	userRepo  UserRepository
	tokenRepo TokenRepository
}

func NewHandler(userRepo *UserRepository, tokenRepo *TokenRepository) *Handler {
	h := &Handler{
		userRepo:  *userRepo,
		tokenRepo: *tokenRepo,
	}
	return h
}

func (h *Handler) Authenticate(login string, password string) (*events.APIGatewayV2HTTPResponse, error) {
	user, err := h.userRepo.FetchUserByLogin(login)

	if err != nil {
		return apiResponse(http.StatusNotFound, ErrorBody{"User not found"})
	}

	if user.IsPasswordValid(password) {
		h.tokenRepo.RegenerageToken(&user)
		user.Password = ""
		return apiResponse(http.StatusOK, user)
	}

	return apiResponse(http.StatusForbidden, ErrorBody{"Wrong login or password"})
}

func (h *Handler) VerifyToken(token string) (*events.APIGatewayV2HTTPResponse, error) {
	userId, err := h.tokenRepo.FetchUserIdByToken(token)

	if err == nil {
		return apiResponse(http.StatusOK, userId)
	}

	return apiResponse(http.StatusNotFound, ErrorBody{"User not authorized"})
}
