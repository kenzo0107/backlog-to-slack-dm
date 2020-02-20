package config

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/kenzo0107/backlog-to-slack-dm/pkg/awsapi"
)

var (
	ssmBacklogAPIKeyName                = os.Getenv("BACKLOG_API_KEY_NAME")
	ssmBacklogBasicAuthEncName          = os.Getenv("BACKLOG_BASIC_AUTH_ENC_NAME")
	ssmSlackBotUserOauthAccessTokenName = os.Getenv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN_NAME")

	ssmParameterKeyNames = []string{
		ssmBacklogAPIKeyName,
		ssmBacklogBasicAuthEncName,
		ssmSlackBotUserOauthAccessTokenName,
	}
)

// Secrets : secrets
var Secrets SecretParameters

// SecretParameters : secret parameters
type SecretParameters struct {
	BacklogAPIKey                string
	BacklogBasicAuthEnc          string
	SlackBotUserOauthAccessToken string
}

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))

	var keys []string
	for _, s := range ssmParameterKeyNames {
		if s == "" {
			continue
		}
		keys = append(keys, s)
	}

	ssmClient := awsapi.NewSSMClient(ssm.New(sess))
	s, err := ssmClient.GetSSMParameters(keys)

	if err != nil {
		log.Fatal(err)
	}

	Secrets.BacklogAPIKey = s[ssmBacklogAPIKeyName]
	Secrets.BacklogBasicAuthEnc = s[ssmBacklogBasicAuthEncName]
	Secrets.SlackBotUserOauthAccessToken = s[ssmSlackBotUserOauthAccessTokenName]
}
