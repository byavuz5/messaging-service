package main

import (
	"sort"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

//Send activity log to activity_logs table.
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
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       []Result          `json:"body"`
	Errors     interface{}       `json:"errors"`
}

type Result struct {
	Username  string `json:"username"`
	Activity  string `json:"activity"`
	CreatedAt string `json:"createdAt"`
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
	result, err := svc.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(request.Username),
			},
		},
		KeyConditionExpression: aws.String("username = :username "),

		TableName: aws.String("activity_logs"),
	})

	// Sorts result from present to past.
	sort.Slice(result.Items, func(i, j int) bool {
		return *result.Items[i]["createdAt"].S > *result.Items[j]["createdAt"].S
	})

	if err != nil {
		err := sendSystemLog("getActivityLogs", err.Error())
		if err != nil {
			response.StatusCode = 400
			response.Errors = err
			return response, nil
		}
		response.StatusCode = 400
		response.Errors = err
		return response, nil
	}
	err = sendActivity(request.Username, "List activity logs.")

	if err != nil {
		err := sendSystemLog("getActivityLogs", err.Error())
		if err != nil {
			response.StatusCode = 400
			response.Errors = err
		}
		response.StatusCode = 400
		response.Errors = err
		return response, nil
	}

	var queryResult []Result

	//Convert dynamodb query output to interface.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &queryResult)
	if err != nil {
		err := sendSystemLog("getActivityLogs", err.Error())
		if err != nil {
			response.StatusCode = 400
			response.Errors = err
			return response, nil
		}
		response.StatusCode = 400
		response.Errors = err
		return response, nil
	}

	response.StatusCode = 200
	response.Body = queryResult

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
