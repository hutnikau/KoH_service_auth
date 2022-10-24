package handlers

import (
	"encoding/json"
	"errors"
	"service-auth/pkg/model"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockedUserRepository struct {
	mock.Mock
}

type MockedTokenRepository struct {
	mock.Mock
}

func (m MockedUserRepository) FetchUserById(userId string) (model.User, error) {
	args := m.Called(userId)
	return args.Get(0).(model.User), args.Error(1)
}
func (m MockedUserRepository) FetchUserByLogin(login string) (model.User, error) {
	args := m.Called(login)
	return args.Get(0).(model.User), args.Error(1)
}

func (m MockedTokenRepository) RegenerateToken(user *model.User) *model.Token {
	args := m.Called(user)
	return args.Get(0).(*model.Token)
}
func (m MockedTokenRepository) FetchUserIdByToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestAuthenticateUserNotFoundFail(test *testing.T) {
	userRepo := MockedUserRepository{}
	tokenRepo := MockedTokenRepository{}
	userRepo.On("FetchUserByLogin", "john").Return(model.User{}, errors.New("could not find user"))

	h := NewHandler(userRepo, tokenRepo)

	response, _ := h.Authenticate("john", "doe")

	if response.StatusCode != 404 {
		test.Errorf("Expected status code 404, %d given", response.StatusCode)
	}
	if response.Body != "{\"error\":\"User not found\"}" {
		test.Errorf("Expected \"User not found\" error message")
	}
}

func TestAuthenticateWrongPasswordFail(test *testing.T) {
	userRepo := MockedUserRepository{}
	tokenRepo := MockedTokenRepository{}
	user := model.User{
		Login:    "john",
		Id:       "123",
		Token:    "foobar",
		Password: "wrong",
	}

	userRepo.On("FetchUserByLogin", "john").Return(user, nil)

	h := NewHandler(userRepo, tokenRepo)

	response, _ := h.Authenticate("john", "doe")

	if response.StatusCode != 403 {
		test.Errorf("Expected status code 403, %d given", response.StatusCode)
	}
	if response.Body != "{\"error\":\"Wrong login or password\"}" {
		test.Errorf("Expected \"Wrong login or password\" error message")
	}
}

func TestAuthenticate(test *testing.T) {
	userRepo := MockedUserRepository{}
	tokenRepo := MockedTokenRepository{}
	user := model.User{
		Login:    "john",
		Id:       "123",
		Token:    "foobar",
		Password: "$2a$10$BK6LS2G/SCVughYKVYJ4i.wOy/uS.rcM8BP.IePjfhqf07jPqi3Zi",
	}
	token := model.Token{
		Token:     "foobar",
		UserId:    "123",
		CreatedAt: time.Now().Unix(),
	}
	userRepo.On("FetchUserByLogin", "john").Return(user, nil)
	tokenRepo.On("RegenerateToken", &user).Return(&token)
	h := NewHandler(userRepo, tokenRepo)
	response, _ := h.Authenticate("john", "doe")
	if response.StatusCode != 200 {
		test.Errorf("Expected status code 200")
	}

	userFromResponse := model.User{}

	_ = json.Unmarshal([]byte(response.Body), &userFromResponse)

	if userFromResponse.Login != "john" || userFromResponse.Id != "123" || userFromResponse.Token != "foobar" {
		test.Errorf("Wrong user data given")
	}

	if userFromResponse.Password != "" {
		test.Errorf("User password must not be in the response")
	}
}
