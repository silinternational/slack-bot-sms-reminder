package nexmo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupEnv() {
	os.Setenv("SLACK_VERIFICATION_TOKEN", "abc")
	os.Setenv("NEXMO_API_KEY", "abc")
	os.Setenv("NEXMO_API_SECRET", "abc")
	os.Setenv("NEXMO_API_FROM", "abc")
}

func TestSendSms(t *testing.T) {
	setupEnv()
	mux := http.NewServeMux()

	// Success response test
	mux.HandleFunc("/sms/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
  "message-count": 1,
  "messages": [
    {
      "to": "447700900000",
      "message-id": "0A0000000123ABCD1",
      "status": "0",
      "remaining-balance": "3.14159265",
      "message-price": "0.03330000",
      "network": "12345"
    }
  ]
}`)
	})

	server := httptest.NewServer(mux)

	err := SendSms(server.URL+"/sms/json", "abc", "abc", "15551212", "15551212", "testing 123")
	if err != nil {
		t.Error(err)
	}

	// Error response
	mux = http.NewServeMux()
	mux.HandleFunc("/sms/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
  "message-count": 1,
  "messages": [
    {
      "status": "2",
      "error-text": "Missing to param"
    }
  ]
}`)
	})

	server = httptest.NewServer(mux)
	err = SendSms(server.URL+"/sms/json", "abc", "abc", "15551212", "15551212", "testing 123")
	if err == nil {
		t.Error("Error not returned when it should have been ")
	}
}
