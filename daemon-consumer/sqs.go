package main

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Receive message fron specified queue.
func Receive(qURL string) error {

    if qURL == "" {
        return errors.New("qURL must not empty")
    }

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(5), // 可視性タイムアウト(秒) 0 to 43,200 (12時間)
		WaitTimeSeconds:     aws.Int64(1), // ロングポーリング(秒) 1 to 20
	})

	if err != nil {
		log.Println("Error", err)
        return err
	}

	if len(result.Messages) == 0 {
		log.Println("Received no messages")
		return nil
	}

	log.Println("Received ", len(result.Messages), " messages")

	for _, message := range result.Messages {

		log.Printf("message(%s) = body: %s, attr: %v",
			*message.MessageId,
			*message.Body,
			message.MessageAttributes)

		_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &qURL,
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			log.Println("Delete Error", err)
			return nil
		}

		messageID := message.MessageId
		log.Printf("Message(%s) Deleted\n", *messageID)
	}

    return nil
}
