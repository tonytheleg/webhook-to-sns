package main

/*
Basic Goals: DONE
Static Payload json created
Send to SNS the payload as message
Test getting message

Next Step Goals: DONE
Use struct for json payload
Publish alert to topic based on provided user

MVP Goals
Send to correct topic based on the request URL
Use STS over user creds?
Determine topic to publish to from calling AWS based on topic owner so topics are not defined in code
*/

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"fmt"
)

type Alert struct {
	AppName        string `json:"AppName"`
	Type           string `json:"Type"`
	AlertID        int64  `json:"AlertID"`
	Summary        string `json:"Summary"`
	OncallEngineer string `json:"oncallEngineer"`
}

func (a *Alert) ToString() string {
	return fmt.Sprintf("\nApp Name: %s\nType: %s\nAlert ID: %d\nSummary: %s\n", a.AppName, a.Type, a.AlertID, a.Summary)
}

func (a *Alert) PublishToTopic(topicArn string) error {
	msg := a.ToString()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-gov-west-1"),
		Credentials: credentials.NewSharedCredentials("", "goalert-test"),
	})
	if err != nil {
		return fmt.Errorf("failed to create a session: %v", err)
	}

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
	sampleTonyAlert := Alert{
		AppName:        "Sample Alert for Tony",
		Type:           "Alert",
		AlertID:        65,
		Summary:        "Job did not complete in time",
		OncallEngineer: "anatale",
	}

	sampleBrianAlert := Alert{
		AppName:        "Sample Alert for Brian",
		Type:           "Alert",
		AlertID:        65,
		Summary:        "Job did not complete in time",
		OncallEngineer: "bsmith",
	}

	tonyTopic := "arn:aws-us-gov:sns:us-gov-west-1:657750906120:goalert_anatale"
	brianTopic := "arn:aws-us-gov:sns:us-gov-west-1:657750906120:goalert_bsmith"

	sampleTonyAlert.PublishToTopic(tonyTopic)
	sampleBrianAlert.PublishToTopic(brianTopic)
}
