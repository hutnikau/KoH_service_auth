package main

import (
	"encoding/json"
	"service-auth/pkg/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type authenticateBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type verifyTokenBody struct {
	Token string `json:"token"`
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	reqJson, _ := json.Marshal(req)
	log.Debug().Msg(string(reqJson))

	log.Debug().RawJSON("jsonfield", reqJson).Msg("Raw json")

	return handlers.UnhandledMethod()

	switch req.Path {
	case "authenticate":
		body := authenticateBody{}
		_ = json.Unmarshal([]byte(req.Body), &body)
		return handlers.Authenticate(body.Login, body.Password)
	case "verify_token":
		body := verifyTokenBody{}
		_ = json.Unmarshal([]byte(req.Body), &body)
		return handlers.VerifyToken(body.Token)
	}
	return handlers.UnhandledMethod()
}
