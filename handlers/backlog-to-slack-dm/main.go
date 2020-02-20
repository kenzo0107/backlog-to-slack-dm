package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kenzo0107/backlog-to-slack-dm/config"
	"github.com/kenzo0107/backlog-to-slack-dm/pkg/utility"
	"github.com/nlopes/slack"

	// fork griffin-stewie/go-backlog because "Add Comment" Activity has no Notification
	backlog "github.com/kenzo0107/go-backlog"
)

var (
	backlogBaseURL = os.Getenv("BACKLOG_BASE_URL")

	// strDeniedToNotifyMails : string of multi email address combined by ","
	// ex. hoge.moge@gmail.com,aiu.eo@gmail.com
	strDeniedToNotifyMails = os.Getenv("DENIED_TO_NOTIFY_MAILS")
	deniedToNotifyMails    []string
)

func init() {
	deniedToNotifyMails = strings.Split(strDeniedToNotifyMails, ",")
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body backlog.Activity

	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		log.Println("json.Unmarshal", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	if len(body.Notifications) == 0 {
		// activities has no notifications
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
		}, nil
	}

	issueURL := fmt.Sprintf("%s/view/%s-%s#comment-%s",
		backlogBaseURL,
		*body.Project.ProjectKey,
		strconv.Itoa(*body.Content.KeyID),
		strconv.Itoa(*body.Content.Comment.ID),
	)

	backlogURL, err := url.Parse(backlogBaseURL)
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{}, err
	}

	backlogClient := backlog.NewClient(
		backlogURL,
		config.Secrets.BacklogAPIKey,
	)
	slackClient := slack.New(config.Secrets.SlackBotUserOauthAccessToken)

	attachment := slack.Attachment{
		Color:     "#42CE9F", // backlog color
		Title:     *body.Content.Summary,
		TitleLink: issueURL,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: fmt.Sprintf("%s's comment", *body.CreatedUser.Name),
				Value: *body.Content.Comment.Content,
				Short: false,
			},
		},
	}

	for _, notification := range body.Notifications {
		// get backlog user information by backlog user ID
		backlogUser, err := backlogClient.User(*notification.User.ID)
		if err != nil {
			log.Println("cannot get backlog user by user ID", err)
			continue
		}

		// Not notified to the email address of the user who refuses notification
		if exists, _, _ := utility.InSlice(*backlogUser.MailAddress, deniedToNotifyMails); exists {
			log.Println("Denied to notify this mail address", *backlogUser.MailAddress)
			continue
		}

		// get slack user by backlog email address
		slackUser, err := slackClient.GetUserByEmail(*backlogUser.MailAddress)
		if err != nil {
			log.Println("cannot get slack user by email", err)
			continue
		}

		// send Slack DM
		_, _, err = slackClient.PostMessage(
			slackUser.ID,
			slack.MsgOptionAttachments(attachment),
		)

		if err != nil {
			log.Println("cannot post Slack DM", *backlogUser.MailAddress, err)
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "OK",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
