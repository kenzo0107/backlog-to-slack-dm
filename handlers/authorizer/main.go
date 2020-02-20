package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kenzo0107/backlog-to-slack-dm/config"
)

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	// get <value> part of "basic <value>" set in header
	a := strings.Split(event.Headers["Authorization"], "Basic ")
	if len(a) <= 1 {
		log.Println("UnAuthorized")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("UnAuthorized")
	}

	enc := a[1]
	if enc != config.Secrets.BacklogBasicAuthEnc {
		log.Println("InvalidToken")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("InvalidToken")
	}

	// authorized and generate policy
	return generatePolicy("user", "Allow", event.MethodArn), nil
}

func main() {
	lambda.Start(handler)
}
