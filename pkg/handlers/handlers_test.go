package handlers

import (
	"encoding/json"
	"service-auth/pkg/model"
	"testing"
)

func TestAuthenticateFail(test *testing.T) {
	userRepo, tokenRepo := NewUserRepository(), NewTokenRepository()

	h := NewHandler(&userRepo, &tokenRepo)
	response, _ := h.Authenticate("a", "b")
	if response.StatusCode != 400 {
		test.Errorf("Expected status code 400")
	}
	if response.Body != "{\"error\":\"User not found\"}" {
		test.Errorf("Expected \"User not found\" error message")
	}
}

func TestAuthenticate(test *testing.T) {
	userRepo, tokenRepo := NewUserRepository(), NewTokenRepository()

	h := NewHandler(&userRepo, &tokenRepo)
	response, _ := h.Authenticate("john", "doe")
	if response.StatusCode != 200 {
		test.Errorf("Expected status code 200")
	}

	user := model.User{}

	_ = json.Unmarshal([]byte(response.Body), &user)

	if user.Login != "john" || user.Id != "123" || user.Token != "foobar" {
		test.Errorf("Wrong user data given")
	}

	if user.Password != "" {
		test.Errorf("User password must not be in the response")
	}
}
