// @@@SNIPSTART go-ipgeo-workflow-imports
package iplocate

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPEND

// @@@SNIPSTART go-ipgeo-workflow-code
// GetAddressFromIP is the Temporal Workflow that retrieves the IP address and location info.
func GetAddressFromIP(ctx workflow.Context, name string) (string, error) {
	// Define the activity options, including the retry policy
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second, //amount of time that must elapse before the first retry occurs
			MaximumInterval:    time.Minute, //maximum interval between retries
			BackoffCoefficient: 2,           //how much the retry interval increases
			// MaximumAttempts: 5, // Uncomment this if you want to limit attempts
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var ipActivities *IPActivities

	var ip string
	err := workflow.ExecuteActivity(ctx, ipActivities.GetIP).Get(ctx, &ip)
	if err != nil {
		return "", fmt.Errorf("Failed to get IP: %s", err)
	}

	var location string
	err = workflow.ExecuteActivity(ctx, ipActivities.GetLocationInfo, ip).Get(ctx, &location)
	if err != nil {
		return "", fmt.Errorf("Failed to get location: %s", err)
	}
	return fmt.Sprintf("Hello, %s. Your IP is %s and your location is %s", name, ip, location), nil
}

// @@@SNIPEND
