package nexmo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const SmsURL = "https://rest.nexmo.com/sms/json"

type SmsRequest struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	From      string `json:"from"`
	To        string `json:"to"`
	Type      string `json:"type"`
	Text      string `json:"text"`
}

type SmsResponse struct {
	MessageCount string `json:"message-count"`
	Messages     []struct {
		To               string `json:"to"`
		MessageID        string `json:"message-id"`
		Status           string `json:"status"`
		RemainingBalance string `json:"remaining-balance"`
		MessagePrice     string `json:"message-price"`
		Network          string `json:"network"`
		ErrorText        string `json:"error-text"`
	} `json:"messages"`
}

func SendSms(apiUrl, apiKey, apiSecret, fromPhone, toPhone, message string) error {

	if apiUrl == "" {
		apiUrl = SmsURL
	}

	body := SmsRequest{
		APIKey:    apiKey,
		APISecret: apiSecret,
		From:      fromPhone,
		To:        toPhone,
		Type:      "text",
		Text:      message,
	}

	js, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json marshal error: %s", err.Error())
	}

	bodyReader := strings.NewReader(string(js))

	req, err := http.NewRequest(http.MethodPost, apiUrl, bodyReader)
	if err != nil {
		return fmt.Errorf("new request error: %s", err.Error())
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %s", err.Error())
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body error: %s", err.Error())
	}

	var smsResponse SmsResponse
	err = json.Unmarshal(respBody, &smsResponse)
	if err != nil {
		return fmt.Errorf("json unmarshal error: %s, body: %s", err.Error(), respBody)
	}

	if smsResponse.Messages[0].Status != "0" {
		return fmt.Errorf("error ending SMS. Code %s, Message: %s", smsResponse.Messages[0].Status, smsResponse.Messages[0].ErrorText)
	}

	return nil
}
