package main

import (
	"encoding/json"
	"service-auth/pkg/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
)

type authenticateBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// reqJson, _ := json.Marshal(req)
	// log.Info().RawJSON("Raw json", reqJson).Msg("Request")

	switch req.RequestContext.RouteKey {
	case "POST /authenticate":
		body := authenticateBody{}
		_ = json.Unmarshal([]byte(req.Body), &body)
		return handlers.Authenticate(body.Login, body.Password)
	case "GET /verify_token":
		if token, ok := req.QueryStringParameters["token"]; ok {
			return handlers.VerifyToken(token)
		}
	}
	return handlers.UnhandledMethod()
}
