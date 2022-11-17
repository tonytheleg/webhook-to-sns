package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"fmt"
)

type Alert struct {
	AppName  string `json:"AppName"`
	Type     string `json:"Type"`
	AlertID  string `json:"AlertID"`
	Summary  string `json:"Summary"`
	TopicArn string `json:"topicArn"`
}

// ToString takes the Alert struct and dumps it into a nice readable format in a SMS text message
func (a *Alert) ToString() string {
	return fmt.Sprintf("\nApp Name: %s\nType: %s\nAlert ID: %s\nSummary: %s\n", a.AppName, a.Type, a.AlertID, a.Summary)
}

// FetchTopic lists the SNS topics in an account and grabs the arn for the topic matching the oncall username
func FetchTopic(sess *session.Session, name string) (string, error) {
	svc := sns.New(sess)

	result, err := svc.ListTopics(nil)
	if err != nil {
		return "", fmt.Errorf("failed to ListTopics: %v", err)
	}

	for _, t := range result.Topics {
		if strings.Contains(*t.TopicArn, name) {
			return *t.TopicArn, nil
		}
	}
	return "", fmt.Errorf("topic not found")
}

// PublishTopic sends the alert to the SNS topic to notify engineers
func (a *Alert) PublishToTopic(sess *session.Session) error {
	msg := a.ToString()
	topicArn := a.TopicArn

	svc := sns.New(sess)
	result, err := svc.Publish(&sns.PublishInput{
		Message:  &msg,
		TopicArn: &topicArn,
	})
	if err != nil {
		return fmt.Errorf("failed to create a session: %v", err)
	}

	log.Println(*result.MessageId)
	return nil
}

func main() {

	// This mimics a potential webhook URL configured in GoAlert for user anatale
	sampleUserURL := "https://goalert.apps.clusterdomain.openshiftusgov.com/anatale"

	// This mimics a sample request with body that the front end may receive with the alert data
	// Its merely here to show how we can parse the request to determine the on-call engineer to page
	requestBody, err := json.Marshal(map[string]string{
		"AppName": "Sample Alert for Tony",
		"Type":    "Alert",
		"AlertID": "65",
		"Summary": "Job did not complete in time",
		"Details": "Sample Alert Output from Alertmanager",
	})
	if err != nil {
		log.Printf("created requestBody failed: %v", err)
	}

	sampleReq, err := http.NewRequest("POST", sampleUserURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("creating sampleReq failed: %v", err)
	}

	// From here on, this is the main logic to parse the alert data and publish to SNS based on info provided

	// Configure AWS session for use later. Instantiang the session with an empty aws.Config
	// will default to using environment variables -- https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
	// You will need AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY exported in your system to work
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-gov-west-1")})
	if err != nil {
		log.Printf("failed to create a session: %v", err)
	}

	// Decodes the request body into an Alert
	var testAlert Alert
	json.NewDecoder(sampleReq.Body).Decode(&testAlert)

	// Determine who is on-call from url
	oncallEng := strings.TrimLeft(sampleReq.URL.Path, "/")

	// Find the topic for on-call !!! Username in Webhook URL in GoAlert must match the username listed in topic !!!
	// Topic format used for test -- goalert_username)
	topicArn, err := FetchTopic(sess, oncallEng)
	if err != nil {
		log.Printf("failed to fetch topic arn %v", err)
	}
	testAlert.TopicArn = topicArn

	// Publish to the topic
	testAlert.PublishToTopic(sess)
}
