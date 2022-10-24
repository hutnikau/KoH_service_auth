package infrastructure

import (
	"errors"
	"service-auth/pkg/handlers"
	"service-auth/pkg/model"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	"github.com/rs/zerolog/log"
)

type DynamodbTokenRepository struct {
	d *dynamodb.DynamoDB
	t string
}

func NewTokenRepository() handlers.TokenRepository {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	d := dynamodb.New(sess)
	r := DynamodbTokenRepository{
		d: d,
		t: "KoH_TokensTable",
	}

	return r
}

func (r DynamodbTokenRepository) RegenerateToken(user *model.User) *model.Token {
	token := model.Token{
		Token:     uuid.New().String(),
		UserId:    (*user).Id,
		CreatedAt: time.Now().Unix(),
	}

	av, err := dynamodbattribute.MarshalMap(token)
	if err != nil {
		log.Fatal().Msgf("Got error marshalling map: %s", err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.t),
	}

	_, err = r.d.PutItem(input)

	if err != nil {
		log.Fatal().Msgf("Error during saving token: %s", err)
	}

	(*user).Token = token.Token
	return &token
}

func (r DynamodbTokenRepository) FetchUserIdByToken(token string) (string, error) {
	tokenItem := model.Token{}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(r.t),
		KeyConditions: map[string]*dynamodb.Condition{
			"token": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(token),
					},
				},
			},
		},
	}

	result, err := r.d.Query(queryInput)
	if err != nil {
		log.Panic().Msgf("Cannot fetch token: %s", err)
	}

	if len(result.Items) == 0 {
		return "", errors.New("could not find token")
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &tokenItem)
	if err != nil {
		log.Panic().Msgf("Failed to unmarshal Token: %s", err)
	}
	return tokenItem.UserId, nil
}
