package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/silinternational/slack-bot-sms-reminder"
	"github.com/silinternational/slack-bot-sms-reminder/lib/db"
	"github.com/silinternational/slack-bot-sms-reminder/lib/nexmo"
	"log"
	"time"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	// Load env config
	err := reminderbot.LoadEnvConfig()
	if err != nil {
		log.Println(err)
		return err
	}

	messages, err := db.ListMessages()
	if err != nil {
		log.Println(err)
		return err
	}

	return ProcessQueuedMessages(messages)
}

func ProcessQueuedMessages(messages []reminderbot.SmsMessage) error {
	now := time.Now().Unix()
	for _, m := range messages {
		if m.SendAt <= now {
			// send message
			err := nexmo.SendSms("", reminderbot.Env.NexmoAPIKey, reminderbot.Env.NexmoAPISecret,
				reminderbot.Env.NexmoAPIFrom, m.PhoneNumber, m.Message)
			if err != nil {
				log.Println("Failed to send SMS, error: ", err.Error())
				continue
			}

			log.Printf("Message ID %s sent\n", m.ID)

			deleted, err := db.DeleteItem(m.ID)
			if err != nil {
				log.Println(err)
				return err
			}
			if deleted {
				log.Printf("Message ID %s deleted\n", m.ID)
			}
		}
	}

	return nil
}
