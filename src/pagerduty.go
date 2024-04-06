// pagerduty.go
package main

import (
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

func getOnCallUser(client *pagerduty.Client) (*pagerduty.User, error) {
	opts := pagerduty.ListOnCallOptions{
		Since: time.Now().Format(time.RFC3339),
		Until: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
	onCalls, err := client.ListOnCalls(opts)
	if err != nil {
		return nil, err
	}
	if len(onCalls.OnCalls) == 0 {
		return nil, fmt.Errorf("no on-call user found")
	}
	// Assuming the first on-call user is the current on-call person
	return &onCalls.OnCalls[0].User, nil
}
