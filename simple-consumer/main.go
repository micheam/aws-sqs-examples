package main

import (
	"log"
	"os"

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
		log.Println("QUEUE_URL is must spesified")
		os.Exit(1)
	}

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
		return
	}

	if len(result.Messages) == 0 {
		log.Println("Received no messages")
		return
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
			return
		}

		messageID := message.MessageId
		log.Printf("Message(%s) Deleted\n", *messageID)
	}
}
