package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/silinternational/slack-bot-sms-reminder"
	"github.com/silinternational/slack-bot-sms-reminder/lib/db"
	"net/http"
	"strings"
	"time"
)

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Load env config
	err := reminderbot.LoadEnvConfig()
	if err != nil {
		return reminderbot.ServerError(err)
	}

	// Convert APIGatewayProxyRequest to *http.Request so SlashCommand can parse it
	httpReq, err := http.NewRequest(req.HTTPMethod, req.Path, strings.NewReader(req.Body))
	if err != nil {
		return reminderbot.ClientError(http.StatusBadRequest, err.Error())
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Parse request into SlashCommand struct
	slashCmd, err := reminderbot.SlashCommandParse(httpReq)
	if err != nil {
		return reminderbot.ClientError(http.StatusBadRequest, err.Error())
	}

	// Authorize request
	if !slashCmd.ValidateToken(reminderbot.Env.VerificationToken) {
		return reminderbot.ClientError(http.StatusForbidden, "Invalid verification token")
	}

	// Check requested command and process
	if slashCmd.Command == reminderbot.CmdSms {
		// Expected input format is 15551212 35m message to send
		parts := strings.Split(slashCmd.Text, " ")
		if len(parts) < 3 {
			return reminderbot.ClientError(
				http.StatusBadRequest,
				"Request must be formatted as \"phone time message\"")
		}

		phone := parts[0]
		sendIn := parts[1]
		msg := strings.Join(parts[2:], " ")

		duration, err := time.ParseDuration(sendIn)
		if err != nil {
			return reminderbot.ClientError(
				http.StatusBadRequest,
				"Unable to parse when to send message, use something like 30m or 1h, got: "+sendIn)
		}

		id := fmt.Sprintf("sms-%s", reminderbot.GetRandString(4))

		sendAt := time.Now().UTC().Add(duration)

		smsMsg := reminderbot.SmsMessage{
			ID:          id,
			SendAt:      sendAt.Unix(),
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			PhoneNumber: phone,
			Message:     msg,
		}

		err = db.PutItem(&smsMsg)
		if err != nil {
			return reminderbot.ServerError(err)
		}

		return reminderbot.ClientError(
			http.StatusOK,
			fmt.Sprintf("SMS reminder scheduled for %s", sendAt.Format(time.RFC3339)))
	}

	return reminderbot.ClientError(http.StatusBadRequest, "Invalid slash command, only support command is /sms")
}
