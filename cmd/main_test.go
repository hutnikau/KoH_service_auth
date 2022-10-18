package main

import (
	"encoding/json"
	"service-auth/pkg/model"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(test *testing.T) {
	request := events.APIGatewayV2HTTPRequest{
		RawPath: "/authenticate",
		Body: `{
			"login" : "john",
			"password" : "doe"
		}`,
	}
	response, _ := handler(request)

	if response.StatusCode != 200 {
		test.Errorf("Expected status code 200")
	}

	user := model.User{}

	_ = json.Unmarshal([]byte(response.Body), &user)

	if user.Login != "john" || user.Id != "123" || user.Token != "foobar" {
		test.Errorf("Wrong user data given")
	}
}
