package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"os"
	"regexp"
)

type Detail struct {
	ScanStatus            string         `json:"scan-status"`
	RepositoryName        string         `json:"repository-name"`
	FindingSeverityCounts map[string]int `json:"finding-severity-counts"`
	ImageDigest           string         `json:"image-digest"`
	ImageTags             []string       `json:"image-tags"`
}

func HandleRequest(ctx context.Context, event events.CloudWatchEvent) error {
	resp := event
	resp.Source = "aws.ecr"
	resp.DetailType = "ECR Image Scan"
	var detail Detail
	json.Unmarshal(event.Detail, &detail)
	if detail.FindingSeverityCounts["CRITICAL"] > 0 || detail.FindingSeverityCounts["HIGH"] > 0 {
		// Find repository name
		var resourceId = regexp.MustCompile(`[^:/]*$`)
		repo := resourceId.FindString(detail.RepositoryName)
		detail.RepositoryName = repo
		detail.ScanStatus = "COMPLETED"
		resp.Detail, _ = json.Marshal(detail)
		message, _ := json.Marshal(resp)
		// Establish a new SNS session
		svc := sns.New(session.New())
		// These are the bare minimum params to send a message.
		params := &sns.PublishInput{
			Message:  aws.String(fmt.Sprintf("%s", message)), // This is the message itself
			TopicArn: aws.String(os.Getenv("sns_arn")),       // This is the ARN of the SNS topic you want to publish to
		}
		_, err := svc.Publish(params)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
