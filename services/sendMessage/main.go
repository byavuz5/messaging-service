package main

import (
	"sort"
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
	Sender  string `json:"sender"`
	Sent_to string `json:"sent_to"`
	Message string `json:"message"`
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
	chatUsers := []string{request.Sender, request.Sent_to}

	//Sort chat users in alphabetical order
	sort.Slice(chatUsers, func(i, j int) bool {
		return chatUsers[i] < chatUsers[j]
	})

	//Create unique room id.
	uniqueRoomID := chatUsers[0] + "-" + chatUsers[1]

	_, err := svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{

			"sender": {
				S: aws.String(request.Sender),
			},

			"createdAt": {
				S: aws.String(time.Now().UTC().String()),
			},
			"message": {
				S: aws.String(request.Message),
			},
			"room_id": {
				S: aws.String(uniqueRoomID),
			},
		},
		TableName: aws.String("messages"),
	})
	if err != nil {
		err := sendSystemLog("sendMessage", err.Error())
		if err != nil {
			response.StatusCode = 400
			response.Errors = err
			return response, nil
		}
		response.StatusCode = 400
		response.Errors = err
		return response, nil

	}
	err = sendActivity(request.Sender, "Send message to "+request.Sent_to+".")
	if err != nil {
		err := sendSystemLog("sendMessage", err.Error())
		if err != nil {
			response.StatusCode = 400
			response.Errors = err
			return response, nil
		}
		response.StatusCode = 400
		response.Errors = err
		return response, nil
	}

	queryResult := Result{
		Message: "Message sended.",
	}

	response.StatusCode = 200
	response.Body = queryResult
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
