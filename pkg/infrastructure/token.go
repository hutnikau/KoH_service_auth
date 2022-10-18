package infrastructure

import (
	"errors"
	"os"
	"service-auth/pkg/model"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/rs/zerolog/log"
)

type TokenRepository struct {
	d *dynamodb.DynamoDB
	t string
}

func NewTokenRepository() TokenRepository {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	d := dynamodb.New(sess)
	r := TokenRepository{
		d: d,
		t: os.Getenv("TOKENS_TABLE_NAME"),
	}

	return r
}

func (r TokenRepository) RegenerageToken(user *model.User) *model.Token {
	token := new(model.Token)
	token.Token = time.Now().String()
	token.UserId = (*user).Id

	av, err := dynamodbattribute.MarshalMap(token)
	if err != nil {
		log.Fatal().Msgf("Got error marshalling map: %s", err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.t),
	}
	log.Info().Msgf("Token %s stored for user %s", token.Token, token.UserId)
	r.d.PutItem(input)

	return token
}

func (r TokenRepository) FetchUserIdByToken(token string) (string, error) {
	result, err := r.d.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(r.t),
		Key: map[string]*dynamodb.AttributeValue{
			"Token": {
				S: aws.String(token),
			},
		},
	})

	if err != nil {
		log.Fatal().Msgf("Cannot fetch token: %s", err)
	}

	if result.Item == nil {
		return "", errors.New("could not find token")
	}

	tokenItem := model.Token{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &tokenItem)
	if err != nil {
		log.Fatal().Msgf("Failed to unmarshal Token: %s", err)
	}
	return tokenItem.UserId, nil
}
