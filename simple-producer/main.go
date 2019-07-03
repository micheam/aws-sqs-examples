package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	// URL to your queue
	qURL := os.Getenv("QUEUE_URL")
    if qURL == "" {
        fmt.Println("QUEUE_URL is must spesified");
        os.Exit(1)
    }

    now := time.Now()
    someID := "i" + string(now.Local().UnixNano())

	result, err := svc.SendMessage(&sqs.SendMessageInput{
        DelaySeconds: aws.Int64(10), // With Timer : 10 seconds
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"originalID": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(someID),
			},
            "processedAt": &sqs.MessageAttributeValue{
                DataType: aws.String("String"),
                StringValue: aws.String(now.Format(time.RFC3339)),
            },
		},
		MessageBody: aws.String("This is sample message."),
		QueueUrl:    &qURL,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", *result.MessageId)
}
