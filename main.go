package main

/*
Basic Goals: DONE
Static Payload json created
Send to SNS the payload as message
Test getting message

Next Step Goals:
API accept requests with json for alert
Publish the alert to a topic

MVP Goals
Send to correct topic based on the request URL
*/

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"fmt"
	"os"
)

// type Alert struct {
// 	AppName string `json:"AppName"`
// 	Type    string `json:"Type"`
// 	AlertID int64  `json:"AlertID`
// 	Summary string `json:"Summary"`
// 	Details string `json:"Details"`
// }

func main() {
	// sample := Alert{
	// 	AppName: "FedRamp GoAlert",
	// 	Type:    "Alert",
	// 	AlertID: 65,
	// 	Summary: "Job did not complete in time",
	// 	Details: "[Prometheus Alertmanager UI](https:///console-openshift-console.apps.bsmith-rosa.fvjw.i1.devshiftusgov.com/monitoring)\n\nJob did not complete in time [View](https://prometheus-k8s-openshift-monitoring.apps.bsmith-rosa.fvjw.i1.devshiftusgov.com/graph?g0.expr=kube_job_spec_completions%7Bjob%3D%22kube-state-metrics%22%2Cnamespace%3D~%22%28openshift-.%2A%7Ckube-.%2A%7Cdefault%29%22%7D+-+kube_job_status_succeeded%7Bjob%3D%22kube-state-metrics%22%2Cnamespace%3D~%22%28openshift-.%2A%7Ckube-.%2A%7Cdefault%29%22%7D+%3E+0&g0.tab=1)\n\n## Payload\n\n```json\n{\n  \"receiver\": \"GoAlert\",\n  \"status\": \"firing\",\n  \"alerts\": [\n    {\n      \"status\": \"firing\",\n      \"labels\": {\n        \"alertname\": \"KubeJobCompletion\",\n        \"container\": \"kube-rbac-proxy-main\",\n        \"endpoint\": \"https-main\",\n        \"job\": \"kube-state-metrics\",\n        \"job_name\": \"ids-tester-27739025\",\n        \"namespace\": \"openshift-suricata\",\n        \"openshift_io_alert_source\": \"platform\",\n        \"prometheus\": \"openshift-monitoring/k8s\",\n        \"service\": \"kube-state-metrics\",\n        \"severity\": \"warning\"\n      },\n      \"annotations\": {\n        \"description\": \"Job openshift-suricata/ids-tester-27739025 is taking more than 12 hours to complete.\",\n        \"summary\": \"Job did not complete in time\"\n      },\n      \"startsAt\": \"2022-09-28T17:05:20.048Z\",\n      \"endsAt\": \"0001-01-01T00:00:00Z\",\n      \"generatorURL\": \"https://prometheus-k8s-openshift-monitoring.apps.bsmith-rosa.fvjw.i1.devshiftusgov.com/graph?g0.expr=kube_job_spec_completions%7Bjob%3D%22kube-state-metrics%22%2Cnamespace%3D~%22%28openshift-.%2A%7Ckube-.%2A%7Cdefault%29%22%7D+-+kube_job_status_succeeded%7Bjob%3D%22kube-state-metrics%22%2Cnamespace%3D~%22%28openshift-.%2A%7Ckube-.%2A%7Cdefault%29%22%7D+%3E+0\\u0026g0.tab=1\",\n      \"fingerprint\": \"19da89970cd63433\"\n    }\n  ],\n  \"groupLabels\": {\n    \"alertname\": \"KubeJobCompletion\",\n    \"severity\": \"warning\"\n  },\n  \"commonLabels\": {\n    \"alertname\": \"KubeJobCompletion\",\n    \"container\": \"kube-rbac-proxy-main\",\n    \"endpoint\": \"https-main\",\n    \"job\": \"kube-state-metrics\",\n    \"job_name\": \"ids-tester-27739025\",\n    \"namespace\": \"openshift-suricata\",\n    \"openshift_io_alert_source\": \"platform\",\n    \"prometheus\": \"openshift-monitoring/k8s\",\n    \"service\": \"kube-state-metrics\",\n    \"severity\": \"warning\"\n  },\n  \"commonAnnotations\": {\n    \"description\": \"Job openshift-suricata/ids-tester-27739025 is taking more than 12 hours to complete.\",\n    \"summary\": \"Job did not complete in time\"\n  },\n  \"externalURL\": \"https:///console-openshift-console.apps.bsmith-rosa.fvjw.i1.devshiftusgov.com/monitoring\",\n  \"version\": \"4\",\n  \"groupKey\": \"{}/{}:{alertname=\\\"KubeJobCompletion\\\", severity=\\\"warning\\\"}\",\n  \"truncatedAlerts\": 0\n}\n\n```",
	// }

	sampleMsg := "{AppName: 'FedRamp GoAlert',Type: 'Alert', AlertID: 65, Summary: 'Job did not complete in time'"
	topic := "arn:aws-us-gov:sns:us-gov-west-1:657750906120:GoAlert-Test-Tony"

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-gov-west-1"),
		Credentials: credentials.NewSharedCredentials("", "goalert-test"),
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := sns.New(sess)

	result, err := svc.Publish(&sns.PublishInput{
		Message:  &sampleMsg,
		TopicArn: &topic,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(*result.MessageId)
}
