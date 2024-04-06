// pagerduty.go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type PagerDutyUser struct {

    ID        string `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
}

func getOnCallUser(client *pagerduty.Client) (*PagerDutyUser, error) {
	ctx:= context.Background()
	opts := pagerduty.ListOnCallOptions{
		Since: time.Now().Format(time.RFC3339),
		Until: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
	onCalls, err := client.ListOnCallsWithContext(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(onCalls.OnCalls) == 0 {
		return nil, fmt.Errorf("no on-call user found")
	}
	// Assuming the first on-call user is the current on-call person
	// Assuming the first on-call user is the current on-call person
    return &PagerDutyUser{
        ID:        onCalls.OnCalls[0].User.ID,
        Name:      onCalls.OnCalls[0].User.Name,
        Email:     onCalls.OnCalls[0].User.Email,
    }, nil
}
