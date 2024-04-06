// main.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/PagerDuty/go-pagerduty"
)

func main() {
	auth := AzureAuth{
		ClientID:       os.Getenv("AZURE_CLIENT_ID"),
		ClientSecret:   os.Getenv("AZURE_CLIENT_SECRET"),
		TenantID:       os.Getenv("AZURE_TENANT_ID"),
		SubscriptionID: os.Getenv("AZURE_SUBSCRIPTION_ID"),
	}

	groupID := os.Getenv("AZURE_AD_GROUP_ID")

	pdToken := os.Getenv("PAGERDUTY_TOKEN")

	pdClient := pagerduty.NewClient(pdToken)

	onCallUser, err := getOnCallUser(pdClient)
	if err != nil {
		log.Fatalf("Failed to get on-call user from PagerDuty: %v", err)
	}

	err = HandleOnCall(auth, groupID, onCallUser)
	if err != nil {
		log.Fatalf("Failed to handle on-call user: %v", err)
	}
}

func HandleOnCall(auth AzureAuth, groupID string, onCallUser *pagerduty.User) error {
	authorizer, err := NewAuthorizerFromClientCredentials(auth)
	if err != nil {
		return err
	}

	graphClient := graphrbac.NewGroupsClient(auth.SubscriptionID)
	graphClient.Authorizer = authorizer

	ctx := context.Background()

	err = addToAzureADGroup(ctx, graphClient, groupID, onCallUser.Email)
	if err != nil {
		return err
	}
	log.Printf("Added user %s to Azure AD group.\n", onCallUser.Name)

	// Simulate the on-call shift ending after 24 hours
	time.Sleep(24 * time.Hour)

	err = removeFromAzureADGroup(ctx, graphClient, groupID, onCallUser.Email)
	if err != nil {
		return err
	}
	log.Printf("Removed user %s from Azure AD group.\n", onCallUser.Name)

	return nil
}
