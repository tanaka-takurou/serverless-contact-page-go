package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type APIResponse struct {
	Message  string `json:"message"`
}

type Response events.APIGatewayProxyResponse

var snsClient *sns.Client

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
						err = sendmessage(ctx, n, b, m)
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

func sendmessage(ctx context.Context, name string, message string, mail string) error {
	if snsClient == nil {
		snsClient = sns.NewFromConfig(getConfig(ctx))
	}

	input := &sns.PublishInput{
		Subject:  aws.String("Serverless Contact"),
		Message:  aws.String("[Name]\n" + name + "\n\n[Mail]\n" + mail + "\n\n[Message]\n" + message),
		TopicArn: aws.String(os.Getenv("TOPIC_ARN")),
	}

	_, err := snsClient.Publish(ctx, input)
	return err
}

func getConfig(ctx context.Context) aws.Config {
	var err error
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Print(err)
	}
	return cfg
}

func main() {
	lambda.Start(HandleRequest)
}
