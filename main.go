package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {

	// here I am Loading the Env variables
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	// Creating a new client to slack by giving token
	// Setting debug to true while developing
	// Also adding a ApplicationToken option to the client
	//One that uses the regular API and one for the websocket events
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	// Socket Mode, this allows the bot to connect via WebSocket.I need to search an alternative, Maybe a Public URL can help
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		// This is the Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	args := os.Args[1:]
	fmt.Println(args)

	preText := "*Hello! My Jenkins build is Finished!*"
	jenkinsURL := "*Build URL:* " + args[0]
	buildResult := "*" + args[1] + "*"
	buildNumber := "*" + args[2] + "*"
	jobName := "*" + args[3] + "*"

	if buildResult == "*SUCCESS*" {
		buildResult = buildResult + " Wohoo🎉"
	} else {
		buildResult = buildResult + ":x:"
	}

	dividerSection1 := slack.NewDividerBlock()
	jenkinsBuildDetails := jobName + " #" + buildNumber + " - " + buildResult + "\n" + jenkinsURL

	preTextField := slack.NewTextBlockObject("mrkdwn", preText+"\n\n", false, false)

	jenkinsBuildDetailsField := slack.NewTextBlockObject("mrkdwn", jenkinsBuildDetails+"\n\n", false, false)

	jenkinsBuildDetailsSection := slack.NewSectionBlock(jenkinsBuildDetailsField, nil, nil)

	preTextSection := slack.NewSectionBlock(preTextField, nil, nil)

	msg := slack.MsgOptionBlocks(
		preTextSection,
		dividerSection1,
		jenkinsBuildDetailsSection,
	)

	_, _, _, err := client.SendMessage(
		os.Getenv("SLACK_CHANNEL_ID"),
		msg,
	)

	if err != nil {
		fmt.Print(err)
		return
	}

	channelID, timestamp, err := client.PostMessage(
		os.Getenv("SLACK_CHANNEL_ID"), slack.MsgOptionText("Hello World!", false),
	)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("Message sent successfully to channel %s at %s", channelID, timestamp)

	// Creating a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())

	// If we fail to cancel the context, the goroutine that WithCancel or WithTimeout created will be retained in memory indefinitely (until the program shuts down), causing a memory leak. If you do this a lot, your memory will balloon significantly. It's best practice to use a defer cancel() immediately after calling WithCancel() or WithTimeout()

	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		// Creating a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// incase context cancel is called we exit the goroutine
			case <-ctx.Done():
				log.Println("🚸Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				// We have a new Events, let's type switch the event

				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("❎Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socketClient.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := handleEventMessage(eventsAPIEvent, client)

					if err != nil {
						// Replacing with actual err handeling
						log.Fatal(err)
					}

				// Handle Slash Commands
				case socketmode.EventTypeSlashCommand:

					command, ok := event.Data.(slack.SlashCommand)
					if !ok {
						log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
						continue
					}
					// handleSlashCommand will take care of the command
					payload, err := handleSlashCommand(command, client)
					if err != nil {
						log.Fatal(err)
					}
					// The payload is the response
					socketClient.Ack(*event.Request, payload)

				case socketmode.EventTypeInteractive:
					interaction, ok := event.Data.(slack.InteractionCallback)
					if !ok {
						log.Printf("Could not type cast the message to a Interaction callback: %v\n", interaction)
						continue
					}

					err := handleInteractionEvent(interaction, client)
					if err != nil {
						log.Fatal(err)
					}
					socketClient.Ack(*event.Request)
					//end of switch

				}

			}
		}
	}(ctx, client, socketClient)

	socketClient.Run()

}

// handleEventMessage will take an event and handle it properly based on the type of event
func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := handleAppMentionEvent(ev, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("❌unsupported event type")
	}
	return nil
}

// handleAppMentionEvent is used to take care of the AppMentionEvent when the bot is mentioned
func handleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {

	// Grabing the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	// Checking if the user said Hello to the bot
	text := strings.ToLower(event.Text)

	// Creating the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Adding Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date⌛",
			Value: time.Now().String(),
		}, {
			Title: "Caller📞",
			Value: user.Name,
		},
	}
	if strings.Contains(text, "hello") {
		// Greet the user
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Pretext = "Greetings💎"
		attachment.Color = "#4af030"
	} else {
		// Send a message to the user
		attachment.Text = fmt.Sprintf("How can I help you %s?", user.Name)
		attachment.Pretext = "Wohoo🎉What is my work today"
		attachment.Color = "#3d3d3d"
	}
	// Sending the message to the channel
	// The Channel  value  is in the event message
	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("📍failed to post message: %w", err)
	}
	return nil
}

// handleSlashCommand will take a slash command and route to the appropriate function
func handleSlashCommand(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	// We need to switch depending on the command
	switch command.Command {
	case "/namaste":
		// This was a hello command, so pass it along to the proper function
		return nil, handleHelloCommand(command, client)
	case "/red-pill-blue-pill":
		return handleChoice(command, client)
	}

	return nil, nil
}

// handleHelloCommand will take care of /namaste submissions
func handleHelloCommand(command slack.SlashCommand, client *slack.Client) error {
	// The Input is found in the text field so
	// Creating the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Adding Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date⌛",
			Value: time.Now().String(),
		}, {
			Title: "Caller📞",
			Value: command.UserName,
		},
	}

	// Greeting the user
	attachment.Text = fmt.Sprintf("Bonjour😏 %s", command.Text)
	attachment.Color = "#4af030"

	// Sending the message to the channel
	// The Channel is available in the command.ChannelID
	_, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("📍failed to post message: %w", err)
	}
	return nil
}

// handleChoice will trigger a Yes or No question to the initializer
func handleChoice(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	// Creating the attachment and assigned based on the message
	attachment := slack.Attachment{}

	// Creating the checkbox element
	checkbox := slack.NewCheckboxGroupsBlockElement("answer",
		slack.NewOptionBlockObject("yes", &slack.TextBlockObject{Text: "Yes", Type: slack.MarkdownType}, &slack.TextBlockObject{Text: "🎡Did you Enjoy it?", Type: slack.MarkdownType}),
		slack.NewOptionBlockObject("no", &slack.TextBlockObject{Text: "No", Type: slack.MarkdownType}, &slack.TextBlockObject{Text: "🎭Did you Dislike it?", Type: slack.MarkdownType}),
	)
	// Creating the Accessory that will be included in the Block and add the checkbox to it
	accessory := slack.NewAccessory(checkbox)
	// Adding Blocks to the attachment
	attachment.Blocks = slack.Blocks{
		BlockSet: []slack.Block{
			// Creating a new section block element and add some text and the accessory to it
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "Did you think I was Helpful🤖?",
				},
				nil,
				accessory,
			),
		},
	}

	attachment.Text = "Rate the experience😊"
	attachment.Color = "#4af030"
	return attachment, nil
}

func handleInteractionEvent(interaction slack.InteractionCallback, client *slack.Client) error {
	// This is where we would handle the interaction
	// Switch depending on the Type
	log.Printf("The action called is: %s\n", interaction.ActionID)
	log.Printf("The response was of type: %s\n", interaction.Type)
	switch interaction.Type {
	case slack.InteractionTypeBlockActions:
		// This is a block action, so we need to handle it

		for _, action := range interaction.ActionCallback.BlockActions {
			log.Printf("%v", action)
			log.Println("Selected option: ", action.SelectedOptions)

		}

	default:

	}

	return nil
}
