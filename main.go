package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	args := os.Args[1:]
	fmt.Println(args)
	godotenv.Load(".env")

	api := slack.New(os.Getenv("SLACK_AUTH_TOKEN"))
	fmt.Printf("%s", os.Getenv("SLACK_AUTH_TOKEN"))
	preText := "*Hello! Your Jenkins build has finished!*"
	jenkinsURL := "*Build URL:* " + args[0]
	buildResult := "*" + args[1] + "*"
	buildNumber := "*" + args[2] + "*"
	jobName := "*" + args[3] + "*"

	if buildResult == "*SUCCESS*" {
		buildResult = buildResult + " :white_check_mark:"
	} else {
		buildResult = buildResult + " :x:"
	}

	dividerSection1 := slack.NewDividerBlock()
	jenkinsBuildDetails := jobName + " #" + buildNumber + " - " + buildResult + "\n" + jenkinsURL
	preTextField := slack.NewTextBlockObject("mrkdwn", preText+"\n\n", false, false)
	jenkinsBuildDetailsField := slack.NewTextBlockObject("mrkdwn", jenkinsBuildDetails, false, false)

	jenkinsBuildDetailsSection := slack.NewSectionBlock(jenkinsBuildDetailsField, nil, nil)
	preTextSection := slack.NewSectionBlock(preTextField, nil, nil)

	msg := slack.MsgOptionBlocks(
		preTextSection,
		dividerSection1,
		jenkinsBuildDetailsSection,
	)

	_, _, _, err := api.SendMessage(
		"C03UKUE4TD2",
		msg,
	)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
