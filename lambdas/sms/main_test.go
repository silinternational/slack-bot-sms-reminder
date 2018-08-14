package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/silinternational/slack-bot-sms-reminder/lib/db"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	db.FlushTables(t)

	body := `token=abc123
&team_id=T0001
&team_domain=example
&enterprise_id=E0001
&enterprise_name=Globular%20Construct%20Inc
&channel_id=C2147483705
&channel_name=test
&user_id=U2147483697
&user_name=Steve
&command=/sms
&text=15551212 35m this is a test
&response_url=https://hooks.slack.com/commands/1234/5678
&trigger_id=13345224609.738474920.8088930838d88f008e0`

	// replace newlines in body
	body = strings.Replace(body, "\n", "", -1)

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/sms",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: body,
	}

	resp, err := handler(req)
	if err != nil {
		t.Error("Unable to create SMS, err: ", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("Did not get 200 back, got: ", resp.StatusCode, " body: ", resp.Body)
	}

	if !strings.Contains(resp.Body, "SMS reminder scheduled for") {
		t.Error("Create SMS body response does not look right, got: ", resp.Body)
	}

	fmt.Fprint(os.Stdout, resp.Body)
}
