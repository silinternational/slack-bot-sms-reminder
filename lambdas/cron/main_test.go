package main

import (
	"github.com/silinternational/slack-bot-sms-reminder"
	"github.com/silinternational/slack-bot-sms-reminder/lib/db"
	"testing"
	"time"
)

func TestProcessQueuedMessages(t *testing.T) {
	db.FlushTables(t)

	now := time.Now().Unix()

	fixtures := []reminderbot.SmsMessage{
		{
			ID:          "sms-sent",
			SendAt:      now - 100,
			CreatedAt:   time.Now().Format(time.RFC3339),
			PhoneNumber: "15551212",
			Message:     "this message will be sent and deleted",
		},
		{
			ID:          "sms-notsent",
			SendAt:      now + 100,
			CreatedAt:   time.Now().Format(time.RFC3339),
			PhoneNumber: "15551212",
			Message:     "this message will not be sent or deleted",
		},
	}

	for _, f := range fixtures {
		err := db.PutItem(f)
		if err != nil {
			t.Error("Unable to put fixture, err: ", err.Error())
		}
	}

	messages, err := db.ListMessages()
	if err != nil {
		t.Error(err)
	}

	if len(messages) != len(fixtures) {
		t.Errorf("Not all fixtures were loaded properly, wanted %v but found %v", len(fixtures), len(messages))
	}

	err = ProcessQueuedMessages(messages)
	if err != nil {
		t.Error(err)
	}

	after, err := db.ListMessages()
	if err != nil {
		t.Error(err)
	}

	if len(after) != 1 {
		t.Error("Message that should have been sent/deleted was not")
	}

}
