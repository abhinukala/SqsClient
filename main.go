package main

import (
	// "context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"sync/atomic"
	"time"
)

func main() {

	var x uint64

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	// Create a new SQS service client
	svc := sqs.New(sess)

	// Set the queue URL
	queueURL := "https://sqs.us-east-1.amazonaws.com/557058020756/TestListner"

	var batchSize = 10 // maximum number of messages to send in a single API call
	var messages []*sqs.SendMessageBatchRequestEntry
	var messageMainQueue []*sqs.SendMessageBatchRequestEntry
	var batchNumberCount int
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				batchNumberCount += 1
				//fmt.Println(messageMainQueue)
				fmt.Println("Total Number of", batchNumberCount)
				fmt.Println("Total Message Sent", len(messageMainQueue))
				
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	for {
		// create a batch of messages to send
		for i := 0; i < batchSize; i++ {
			atomic.AddUint64(&x, 1)
			message := &sqs.SendMessageBatchRequestEntry{
				Id:          aws.String(fmt.Sprintf("message-%d", x)),
				MessageBody: aws.String(fmt.Sprintf("Test Message-%d", x)),
			}
			messages = append(messages, message)
			messageMainQueue = append(messageMainQueue, message)

		}

		// send the batch of messages
		_, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
			QueueUrl: aws.String(queueURL),
			Entries:  messages,
		})
		if err != nil {
			// handle error
			fmt.Println(err)
			fmt.Println("error")
		}

		// reset the messages slice
		messages = nil
		
		
		time.Sleep(time.Second/20 )
		// wait for the required amount of time before sending the next batch of messages
	}
	
}
