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
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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

func (a *Alert) ToString() string {
	return fmt.Sprintf("\nApp Name: %s\nType: %s\nAlert ID: %s\nSummary: %s\n", a.AppName, a.Type, a.AlertID, a.Summary)
}

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

	/*
		This test assumes that the front half of the webhook will:
		* Accept a request from GoAlert with the alert data
		* The request data is parseable, specifically the Request URL
		* The webhook URL for each user is unique and contains a single path with users username

		With that data we can:
		* Parse the user from the URL thats on-call
		* Fetch the topic for the user
		* Send the alert data from GoAlert to the topic and notify the engineer
	*/

	// This mimics a potential webhook URL configured in GoAlert for user anatale
	sampleUserURL := "https://goalert.apps.appsrefrs01ugw1.p1.openshiftusgov.com/anatale"

	// This mimics a sample request with body that the front end may receive with the alert data
	// Its merely here to show how we can parse the request to determine the on-call engineer to page
	requestBody, err := json.Marshal(map[string]string{
		"AppName": "Sample Alert for Tony",
		"Type":    "Alert",
		"AlertID": "65",
		"Summary": "Job did not complete in time",
		"Details": "[Prometheus Alertmanager UI](https:///console-openshift-console.apps.bsmith-rosa.fvjw.i1.devshiftusgov.com/monitoring)\n\nJob did not complete in time [View](https://prometheus-k8s-openshift-monitoring.apps.bsmith-rosa.fvjw.i1.devshiftusgov.com/graph?g0.expr=kube_job_spec_completions%7Bjob%3D%22kube-state-metrics%22%2Cnamespace%3D~%22%28openshift-.%2A%7Ckube-.%2A%7Cdefault%29%22%7D+-+kube_job_status_succeeded%7Bjob%3D%22kube-state-metrics%22%2Cnamespace%3D~%22%28openshift-.%2A%7Ckube-.%2A%7Cdefault%29%22%7D+%3E+0&g0.tab=1)\n\n## Payload\n\n```json\n{\n  \"receiver\": \"GoAlert\",\n  \"status\": \"firing\",\n  \"alerts\": [\n    {\n      \"status\": \"firing\",\n      \"labels\": {\n        \"alertname\": \"KubeJobCompletion\",\n        \"container\": \"kube-rbac-proxy-main\",\n        \"endpoint\": \"https-main\",\n        \"job\": \"kube-state-metrics\",\n        \"job_name\": \"ids-tester-27739025\",\n        \"namespace\": \"openshift-suricata\",\n        \"openshift_io_alert_source\": \"platform\",\n        \"prometheus\": \"openshift-monitoring/k8s\",\n        \"service\": \"kube-state-metrics\",\n        \"severity\": \"warning\"\n      },\n      \"annotations\": {\n        \"description\": \"Job openshift-suricata/ids-tester-27739025 is taking more than 12 hours to complete.\",\n        \"summary\": \"Job did not complete in time\"\n      },\n      \"startsAt\": \"2022-09-28T17:05:20.048Z\",\n      \"endsAt\": \"0001-01-01T00:00:00Z\",\n      \"generatorURL\": \"https://prometheus-k8s-openshift-monitoring.apps.bsmith-rosa.fvjw.i1.devshiftusgov.com/graph?g0.expr=kube_job_spec_completions%7Bjob%3D%22kube-state-metrics%22%2Cnamespace%3D~%22%28openshift-.%2A%7Ckube-.%2A%7Cdefault%29%22%7D+-+kube_job_status_succeeded%7Bjob%3D%22kube-state-metrics%22%2Cnamespace%3D~%22%28openshift-.%2A%7Ckube-.%2A%7Cdefault%29%22%7D+%3E+0\\u0026g0.tab=1\",\n      \"fingerprint\": \"19da89970cd63433\"\n    }\n  ],\n  \"groupLabels\": {\n    \"alertname\": \"KubeJobCompletion\",\n    \"severity\": \"warning\"\n  },\n  \"commonLabels\": {\n    \"alertname\": \"KubeJobCompletion\",\n    \"container\": \"kube-rbac-proxy-main\",\n    \"endpoint\": \"https-main\",\n    \"job\": \"kube-state-metrics\",\n    \"job_name\": \"ids-tester-27739025\",\n    \"namespace\": \"openshift-suricata\",\n    \"openshift_io_alert_source\": \"platform\",\n    \"prometheus\": \"openshift-monitoring/k8s\",\n    \"service\": \"kube-state-metrics\",\n    \"severity\": \"warning\"\n  },\n  \"commonAnnotations\": {\n    \"description\": \"Job openshift-suricata/ids-tester-27739025 is taking more than 12 hours to complete.\",\n    \"summary\": \"Job did not complete in time\"\n  },\n  \"externalURL\": \"https:///console-openshift-console.apps.bsmith-rosa.fvjw.i1.devshiftusgov.com/monitoring\",\n  \"version\": \"4\",\n  \"groupKey\": \"{}/{}:{alertname=\\\"KubeJobCompletion\\\", severity=\\\"warning\\\"}\",\n  \"truncatedAlerts\": 0\n}\n\n```",
	})
	if err != nil {
		log.Printf("created requestBody failed: %v", err)
	}

	sampleReq, err := http.NewRequest("POST", sampleUserURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("creating sampleReq failed: %v", err)
	}

	// From here on, this is parsing the data and publishing to SNS based on info provided
	// Decodes the request body into an Alert
	var testAlert Alert
	json.NewDecoder(sampleReq.Body).Decode(&testAlert)

	// Determine who is on-call from url
	oncallEng := strings.TrimLeft(sampleReq.URL.Path, "/")

	// Find the topic for on-call !!!Username in Webhook URL in GoAlert must match the username listed in topic!!!
	// Topic format (temporarily for test -- goalert_username)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-gov-west-1"),
		Credentials: credentials.NewSharedCredentials("", "goalert-test"),
	})
	if err != nil {
		log.Printf("failed to create a session: %v", err)
	}

	topicArn, err := FetchTopic(sess, oncallEng)
	if err != nil {
		log.Printf("failed to fetch topic arn %v", err)
	}
	testAlert.TopicArn = topicArn

	// Publish to the topic
	testAlert.PublishToTopic(sess)
}
