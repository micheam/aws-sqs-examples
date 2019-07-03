package main

import (
	"errors"
	"time"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Send message into specified queue.
func Send(qURL string) error {

    if qURL == "" {
        return errors.New("qURL must not empty")
    }

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

    now := time.Now()
    someID := "i" + strconv.FormatInt(now.Local().UnixNano(), 10)

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
		return err
	}

	log.Println("Send Success", *result.MessageId)
    return nil
}
