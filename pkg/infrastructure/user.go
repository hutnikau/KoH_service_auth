package infrastructure

import (
	"errors"
	"service-auth/pkg/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/rs/zerolog/log"
)

type UserRepository struct {
	d *dynamodb.DynamoDB
	t string
}

func NewUserRepository() UserRepository {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	d := dynamodb.New(sess)
	r := UserRepository{
		d: d,
		t: "KoH_UsersTable",
	}

	return r
}

func (userRepo UserRepository) FetchUserById(userId string) (model.User, error) {
	result, err := userRepo.d.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(userRepo.t),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(userId),
			},
		},
	})

	if err != nil {
		log.Fatal().Msgf("Cannot fetch user: %s", err)
	}

	userItem := model.User{}

	if result.Item == nil {
		return userItem, errors.New("could not find user by id")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &userItem)
	if err != nil {
		log.Fatal().Msgf("Failed to unmarshal Token: %s", err)
	}
	return userItem, nil
}

func (userRepo UserRepository) FetchUserByLogin(login string) (model.User, error) {
	userItem := model.User{}

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(userRepo.t),
		IndexName: aws.String("KoH_UsersLogin"),
		KeyConditions: map[string]*dynamodb.Condition{
			"login": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(login),
					},
				},
			},
		},
	}

	result, err := userRepo.d.Query(queryInput)

	if err != nil {
		log.Panic().Msgf("Cannot fetch user by login: %s", err)
	}

	if len(result.Items) == 0 {
		return userItem, errors.New("could not find user by login")
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &userItem)
	if err != nil {
		log.Panic().Msgf("Failed to unmarshal User: %s", err)
	}
	return userItem, nil
}
