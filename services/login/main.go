package main

import (
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

//Send user activity log to activity_logs table.

func sendActivity(username string, activity string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	_, err := svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{

			"username": {
				S: aws.String(username),
			},

			"createdAt": {
				S: aws.String(time.Now().UTC().String()),
			},
			"activity": {
				S: aws.String(activity),
			},
		},
		TableName: aws.String("activity_logs"),
	})
	if err != nil {

		err := sendSystemLog("sendActivityLog", err.Error())
		if err != nil {
			return err
		}
		return err
	}
	return nil
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

	result, err := svc.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(request.Username),
			},
		},
		KeyConditionExpression: aws.String("username = :username "),

		TableName: aws.String("users"),
	})

	if err != nil {
		err := sendSystemLog("login", err.Error())
		if err != nil {
			response.StatusCode = 400
			response.Errors = err
			return response, nil
		}
		response.StatusCode = 400
		response.Errors = err
		return response, nil
	}
	//Check result length to check if the user has created a record before.
	if len(result.Items) > 0 {
		//Check if the user's password is correct.
		if *result.Items[0]["password"].S == request.Password {
			err := sendActivity(request.Username, "Success login.")
			if err != nil {
				err := sendSystemLog("login", err.Error())
				if err != nil {
					response.StatusCode = 400
					response.Errors = err
					return response, nil
				}
				response.StatusCode = 400
				response.Errors = err
				return response, nil
			}

			queryResult.Message = "Success login, username: " + request.Username

		} else {
			err := sendActivity(request.Username, "Invalid login.")
			if err != nil {
				err := sendSystemLog("login", err.Error())
				if err != nil {
					response.StatusCode = 400
					response.Errors = err
					return response, nil
				}
				response.StatusCode = 400
				response.Errors = err
				return response, nil
			}
			queryResult.Message = "Wrong password."
		}
	} else {
		queryResult.Message = "The username you entered doesn't belong to an account."

	}
	response.StatusCode = 200
	response.Body = queryResult

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
