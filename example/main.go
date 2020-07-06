package main

import (
	"flag"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/kent-hamaguchi/aws-cost-reporter/reporter"
)

func main() {
	prof := flag.String("prof", "", "AWS Prpfile name")
	webhook := flag.String("webhook", "", "Slack webhook URL")
	flag.Parse()
	cfg, err := external.LoadDefaultAWSConfig(external.WithSharedConfigProfile(*prof))
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	r := reporter.New(reporter.Config{
		AWSConfig:       cfg,
		SlackWebhookURL: *webhook,
	})
	if err := r.Send(); err != nil {
		panic("report error, " + err.Error())
	}
}
