// @@@SNIPSTART go-ipgeo-workflow-test-setup
package iplocate_test

import (
	"testing"

	"temporal-ip-geolocation/iplocate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go.temporal.io/sdk/testsuite"
)

// @@@SNIPEND

// @@@SNIPSTART go-ipgeo-workflow-test-workflow

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	activities := &iplocate.IPActivities{}

	// Mock activity implementation
	env.OnActivity(activities.GetIP, mock.Anything).Return("1.1.1.1", nil)
	env.OnActivity(activities.GetLocationInfo, mock.Anything, "1.1.1.1").Return("Planet Earth", nil)

	env.ExecuteWorkflow(iplocate.GetAddressFromIP, "Temporal")

	var result string
	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, "Hello, Temporal. Your IP is 1.1.1.1 and your location is Planet Earth", result)
}

// @@@SNIPEND
