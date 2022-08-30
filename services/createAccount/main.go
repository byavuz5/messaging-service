package main

import (
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//Send system log to system_logs table.

func sendSystemLog(service_name string, err_message string) error {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	_, err := svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{

			"service_name": {
				S: aws.String(service_name),
			},

			"createdAt": {
				S: aws.String(time.Now().UTC().String()),
			},
			"err_message": {
				S: aws.String(err_message),
			},
		},
		TableName: aws.String("system_logs"),
	})

	if err != nil {
		return err
	} else {
		return nil
	}
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       Result            `json:"body"`
	Errors     interface{}       `json:"errors"`
}

type Result struct {
	Message string `json:"message"`
}

func Handler(request Request) (Response, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	response := Response{
		Headers: map[string]string{"Content-Type": "application/json"},
	}

	queryResult := Result{}

	_, err := svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{

			"username": {
				S: aws.String(request.Username),
			},
			"password": {
				S: aws.String(request.Password),
			},
			"createdAt": {
				S: aws.String(time.Now().UTC().String()),
			},
		},
		TableName: aws.String("users"),
		//Check username exists on the users table.
		ConditionExpression: aws.String("attribute_not_exists(username)"),
	})

	if err != nil {
		// Check the error code to if the username exists from the table.
		res := strings.Contains(err.Error(), "ConditionalCheckFailedException")
		if res {
			queryResult.Message = "Username already exist."
			response.StatusCode = 200
			response.Body = queryResult
			return response, nil
		} else {

			err := sendSystemLog("createAccount", err.Error())
			if err != nil {
				response.StatusCode = 400
				response.Errors = err
				return response, nil
			}
			response.StatusCode = 400
			response.Errors = err
			return response, nil
		}
	}

	queryResult.Message = "User created, username: " + request.Username

	response.StatusCode = 200
	response.Body = queryResult
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
