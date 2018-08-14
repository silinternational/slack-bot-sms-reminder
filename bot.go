package reminderbot

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const CmdSms = "/sms"
const DynamoDBTableName = "reminderbot_queue"

var Client *slack.Client
var Env envConfig

// Log errors to stderr
var ErrorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type envConfig struct {
	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"SLACK_VERIFICATION_TOKEN" required:"true"`

	// Nexmo credentials for sending SMS
	NexmoAPIKey    string `envconfig:"NEXMO_API_KEY" required:"true"`
	NexmoAPISecret string `envconfig:"NEXMO_API_SECRET" required:"true"`
	NexmoAPIFrom   string `envconfig:"NEXMO_API_FROM" required:"true"`
}

type SmsMessage struct {
	ID          string
	SendAt      int64
	CreatedAt   string
	PhoneNumber string
	Message     string
}

// SlashCommand contains information about a request of the slash command
type SlashCommand struct {
	Token          string `json:"token"`
	TeamID         string `json:"team_id"`
	TeamDomain     string `json:"team_domain"`
	EnterpriseID   string `json:"enterprise_id,omitempty"`
	EnterpriseName string `json:"enterprise_name,omitempty"`
	ChannelID      string `json:"channel_id"`
	ChannelName    string `json:"channel_name"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	Command        string `json:"command"`
	Text           string `json:"text"`
	ResponseURL    string `json:"response_url"`
	TriggerID      string `json:"trigger_id"`
}

// SlashCommandParse will parse the request of the slash command
func SlashCommandParse(r *http.Request) (s SlashCommand, err error) {
	if err = r.ParseForm(); err != nil {
		return s, err
	}
	s.Token = r.PostForm.Get("token")
	s.TeamID = r.PostForm.Get("team_id")
	s.TeamDomain = r.PostForm.Get("team_domain")
	s.EnterpriseID = r.PostForm.Get("enterprise_id")
	s.EnterpriseName = r.PostForm.Get("enterprise_name")
	s.ChannelID = r.PostForm.Get("channel_id")
	s.ChannelName = r.PostForm.Get("channel_name")
	s.UserID = r.PostForm.Get("user_id")
	s.UserName = r.PostForm.Get("user_name")
	s.Command = r.PostForm.Get("command")
	s.Text = r.PostForm.Get("text")
	s.ResponseURL = r.PostForm.Get("response_url")
	s.TriggerID = r.PostForm.Get("trigger_id")
	return s, nil
}

// ValidateToken validates verificationTokens
func (s SlashCommand) ValidateToken(verificationTokens ...string) bool {
	for _, token := range verificationTokens {
		if s.Token == token {
			return true
		}
	}
	return false
}

func LoadEnvConfig() error {
	err := envconfig.Process("", &Env)
	if err != nil {
		return err
	}

	return nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func ServerError(err error) (events.APIGatewayProxyResponse, error) {
	ErrorLogger.Println(err.Error())
	js, _ := json.Marshal(http.StatusText(http.StatusInternalServerError))
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(js),
	}, err
}

// Similarly add a helper for send responses relating to client errors.
func ClientError(status int, body string) (events.APIGatewayProxyResponse, error) {
	js, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(js),
	}, nil
}

// GetTableName returns the env var value of the string passed in or the string itself
func GetDbTableName(table string) string {
	envOverride := os.Getenv(table)
	if envOverride != "" {
		return envOverride
	}

	return table
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GetRandString returns a random string of given length
func GetRandString(length int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
