// azure.go
package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

type AzureAuth struct {
	ClientID       string
	ClientSecret   string
	TenantID       string
	SubscriptionID string
}

func NewAuthorizerFromClientCredentials(auth AzureAuth) (*autorest.Authorizer, error) {
	oauthConfig, err := adal.NewOAuthConfig("https://login.microsoftonline.com/"+auth.TenantID, auth.ClientID)
	if err != nil {
		return nil, err
	}

	spt, err := adal.NewServicePrincipalToken(*oauthConfig, auth.ClientID, auth.ClientSecret, "https://graph.microsoft.com/")
	if err != nil {
		return nil, err
	}

	return autorest.NewBearerAuthorizer(spt), nil
}

func addToAzureADGroup(ctx context.Context, client graphrbac.GroupsClient, groupID, userEmail string) error {
	_, err := client.AddMember(ctx, groupID, graphrbac.GroupAddMemberParameters{
		URL:          fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userEmail),
		ODataBind:    to.StringPtr(fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userEmail)),
		ODataContext: to.StringPtr("https://graph.microsoft.com/v1.0/$metadata#groups('group_id')/members/$entity"),
	})
	if err != nil {
		return err
	}
	return nil
}

func removeFromAzureADGroup(ctx context.Context, client graphrbac.GroupsClient, groupID, userEmail string) error {
	_, err := client.RemoveMember(ctx, groupID, graphrbac.GroupRemoveMemberParameters{
		URL:          fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userEmail),
		ODataBind:    to.StringPtr(fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userEmail)),
		ODataContext: to.StringPtr("https://graph.microsoft.com/v1.0/$metadata#groups('group_id')/members/$entity"),
	})
	if err != nil {
		return err
	}
	return nil
}
