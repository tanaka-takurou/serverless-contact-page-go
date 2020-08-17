package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIResponse struct {
	Message  string `json:"message"`
}

type Response events.APIGatewayProxyResponse

const layout   string = "2006-01-02 15:04"

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "sendmessage" :
			if n, ok := d["name"]; ok {
				if b, ok := d["message"]; ok {
					if m, ok := d["mail"]; ok {
						err = sendmessage(n, b, m)
					}
				}
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		log.Print(err)
		jsonBytes, _ = json.Marshal(APIResponse{Message: fmt.Sprint(err)})
		return Response{
			StatusCode: 500,
			Body: string(jsonBytes),
		}, nil
	}
	jsonBytes, _ = json.Marshal(APIResponse{Message: "Success"})
	return Response {
		StatusCode: 200,
		Body: string(jsonBytes),
	}, nil
}

func sendmessage(name string, message string, mail string) error {
	svc := sns.New(session.New(), &aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})

	input := &sns.PublishInput{
		Message:  aws.String("[Name]\n" + name + "\n\n[Mail]\n" + mail + "\n\n[Message]\n" + message),
		TopicArn: aws.String(os.Getenv("TOPIC_ARN")),
	}
	_, err := svc.Publish(input)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
