package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func sendSystemLog(json string) {
	url := os.Getenv("API_URL") + "sendSystemLog"
	fmt.Println("URL:>", url)
	fmt.Println(json)
	var jsonStr = []byte(json)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

type Request struct {
	Username string `json:"username"`
	Activity string `json:"activity"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(request Request) (Response, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	response := Response{
		Message: "",
		Ok:      true,
	}
	_, err := svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{

			"username": {
				S: aws.String(request.Username),
			},

			"createdAt": {
				S: aws.String(time.Now().UTC().String()),
			},
			"activity": {
				S: aws.String(request.Activity),
			},
		},
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
	})
	if err != nil {

		sendSystemLog("{\"service_name\":\"sendActivityLog\",\"err_message\":\"" + err.Error() + "\"}")
		response.Message = err.Error()
		return response, nil
	}

	response.Message = "Activity sended."
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
