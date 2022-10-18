package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var ErrorMethodNotAllowed = "method Not allowed"

func apiResponse(status int, body interface{}) (*events.APIGatewayV2HTTPResponse, error) {
	resp := events.APIGatewayV2HTTPResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

func UnhandledMethod() (*events.APIGatewayV2HTTPResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
