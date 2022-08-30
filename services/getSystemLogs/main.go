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

type Request struct {
	Service_name string `json:"service_name"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       []Result          `json:"body"`
	Errors     interface{}       `json:"errors"`
}

type Result struct {
	Service_name string `json:"service_name"`
	Err_message  string `json:"err_message"`
	CreatedAt    string `json:"createdAt"`
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
			":service_name": {
				S: aws.String(request.Service_name),
			},
		},
		KeyConditionExpression: aws.String("service_name = :service_name "),

		TableName: aws.String("system_logs"),
	})
	// Sorts result from present to past.
	sort.Slice(result.Items, func(i, j int) bool {
		return *result.Items[i]["createdAt"].S > *result.Items[j]["createdAt"].S
	})
	if err != nil {
		err := sendSystemLog("getSystemLogs", err.Error())
		if err != nil {
			response.StatusCode = 400
			response.Errors = err
			return response, nil
		}
		response.StatusCode = 400
		response.Errors = err
		return response, nil
	}

	var queryResult []Result

	//Convert dynamodb query output to interface.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &queryResult)
	if err != nil {
		err := sendSystemLog("getSystemLogs", err.Error())
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
