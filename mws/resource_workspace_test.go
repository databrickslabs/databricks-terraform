package mws

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/databrickslabs/terraform-provider-databricks/common"

	"github.com/databrickslabs/terraform-provider-databricks/qa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMwsAccWorkspace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	acctID := qa.GetEnvOrSkipTest(t, "DATABRICKS_ACCOUNT_ID")
	client := common.CommonEnvironmentClient()
	workspaceList, err := NewWorkspacesAPI(context.Background(), client).List(acctID)
	assert.NoError(t, err, err)
	t.Log(workspaceList)
}

func TestResourceWorkspaceCreate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/accounts/abc/workspaces",
				ExpectedRequest: Workspace{
					AccountID:                           "abc",
					IsNoPublicIPEnabled:                 true,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					StorageCustomerManagedKeyID:         "def",
				},
				Response: Workspace{
					WorkspaceID:    1234,
					AccountID:      "abc",
					DeploymentName: "900150983cd24fb0",
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/accounts/abc/workspaces/1234",
				Response: Workspace{
					WorkspaceID:                         1234,
					WorkspaceStatus:                     WorkspaceStatusRunning,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					StorageCustomerManagedKeyID:         "def",
					AccountID:                           "abc",
				},
			},
		},
		Resource: ResourceWorkspace(),
		State: map[string]interface{}{
			"account_id":     "abc",
			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
		},
		Create: true,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc/1234", d.Id())
}

func TestResourceWorkspaceCreateWithIsNoPublicIPEnabledFalse(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/accounts/abc/workspaces",
				ExpectedRequest: Workspace{
					AccountID:                           "abc",
					IsNoPublicIPEnabled:                 false,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					StorageCustomerManagedKeyID:         "def",
				},
				Response: Workspace{
					WorkspaceID:    1234,
					AccountID:      "abc",
					DeploymentName: "900150983cd24fb0",
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/accounts/abc/workspaces/1234",
				Response: Workspace{
					WorkspaceID:                         1234,
					WorkspaceStatus:                     WorkspaceStatusRunning,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					StorageCustomerManagedKeyID:         "def",
					AccountID:                           "abc",
				},
			},
		},
		Resource: ResourceWorkspace(),
		State: map[string]interface{}{
			"account_id":     "abc",
			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"is_no_public_ip_enabled":                  false,
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
		},
		Create: true,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc/1234", d.Id())
}

func TestResourceWorkspaceCreateLegacyConfig(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/accounts/abc/workspaces",
				ExpectedRequest: Workspace{
					AccountID:                           "abc",
					IsNoPublicIPEnabled:                 true,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
				},
				Response: Workspace{
					WorkspaceID:    1234,
					AccountID:      "abc",
					DeploymentName: "900150983cd24fb0",
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/accounts/abc/workspaces/1234",
				Response: Workspace{
					WorkspaceID:                         1234,
					WorkspaceStatus:                     WorkspaceStatusRunning,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					AccountID:                           "abc",
				},
			},
		},
		Resource: ResourceWorkspace(),
		State: map[string]interface{}{
			"account_id":               "abc",
			"aws_region":               "us-east-1",
			"credentials_id":           "bcd",
			"customer_managed_key_id":  "def",
			"deployment_name":          "900150983cd24fb0",
			"workspace_name":           "labdata",
			"network_id":               "fgh",
			"storage_configuration_id": "ghi",
		},
		Create: true,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc/1234", d.Id())
}

func TestResourceWorkspaceCreate_Error(t *testing.T) {
	t.Skipf("Making this test skip until we can configure sleep timings for test purposes")
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/accounts/abc/workspaces",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/accounts/abc/workspaces",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceWorkspace(),
		State: map[string]interface{}{
			"account_id":     "abc",
			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"is_no_public_ip_enabled":                  true,
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
		},
		Create: true,
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "", d.Id(), "Id should be empty for error creates")
}

func TestResourceWorkspaceRead(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/accounts/abc/workspaces/1234",
				Response: Workspace{
					AccountID:                           "abc",
					WorkspaceStatus:                     WorkspaceStatusRunning,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					StorageCustomerManagedKeyID:         "def",
					WorkspaceID:                         1234,
				},
			},
		},
		Resource: ResourceWorkspace(),
		Read:     true,
		New:      true,
		ID:       "abc/1234",
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc/1234", d.Id(), "Id should not be empty")
	assert.Equal(t, "us-east-1", d.Get("aws_region"))
	assert.Equal(t, "bcd", d.Get("credentials_id"))
	assert.Equal(t, "def", d.Get("managed_services_customer_managed_key_id"))
	assert.Equal(t, "def", d.Get("storage_customer_managed_key_id"))
	assert.Equal(t, "900150983cd24fb0", d.Get("deployment_name"))
	assert.Equal(t, true, d.Get("is_no_public_ip_enabled"))
	assert.Equal(t, "fgh", d.Get("network_id"))
	assert.Equal(t, "ghi", d.Get("storage_configuration_id"))
	assert.Equal(t, 1234, d.Get("workspace_id"))
	assert.Equal(t, "labdata", d.Get("workspace_name"))
	assert.Equal(t, "RUNNING", d.Get("workspace_status"))
	assert.Equal(t, "https://900150983cd24fb0.cloud.databricks.com", d.Get("workspace_url"))
}

func TestResourceWorkspaceRead_Issue382(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/accounts/abc/workspaces/1234",
				Response: Workspace{
					AccountID:                           "abc",
					WorkspaceStatus:                     WorkspaceStatusRunning,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "prefix-900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					StorageCustomerManagedKeyID:         "def",
					WorkspaceID:                         1234,
				},
			},
		},
		InstanceState: map[string]string{
			"account_id":     "abc",
			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
		},
		State: map[string]interface{}{
			"account_id":     "abc",
			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
		},
		Resource: ResourceWorkspace(),
		Read:     true,
		New:      true,
		ID:       "abc/1234",
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc/1234", d.Id(), "Id should not be empty")
	assert.Equal(t, "prefix-900150983cd24fb0", d.Get("deployment_name"))
	assert.Equal(t, "https://prefix-900150983cd24fb0.cloud.databricks.com", d.Get("workspace_url"))
}

func TestResourceWorkspaceRead_NotFound(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
				Response: common.APIErrorBody{
					ErrorCode: "NOT_FOUND",
					Message:   "Item not found",
				},
				Status: 404,
			},
		},
		Resource: ResourceWorkspace(),
		Read:     true,
		Removed:  true,
		ID:       "abc/1234",
	}.ApplyNoError(t)
}

func TestResourceWorkspaceRead_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{ // read log output for correct url...
				Method:   "GET",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceWorkspace(),
		Read:     true,
		ID:       "abc/1234",
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "abc/1234", d.Id(), "Id should not be empty for error reads")
}

func TestResourceWorkspaceUpdate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "PATCH",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
				ExpectedRequest: map[string]interface{}{
					"credentials_id":                  "bcd",
					"network_id":                      "fgh",
					"storage_customer_managed_key_id": "def",
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/accounts/abc/workspaces/1234",
				Response: Workspace{
					WorkspaceStatus:                     WorkspaceStatusRunning,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					StorageCustomerManagedKeyID:         "def",
					AccountID:                           "abc",
					WorkspaceID:                         1234,
				},
			},
		},
		Resource: ResourceWorkspace(),
		State: map[string]interface{}{
			"account_id":     "abc",
			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"is_no_public_ip_enabled":                  true,
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
			"workspace_id":                             1234,
		},
		Update: true,
		ID:     "abc/1234",
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc/1234", d.Id(), "Id should be the same as in reading")
}

func TestResourceWorkspaceUpdate_NotAllowed(t *testing.T) {
	qa.ResourceFixture{
		Resource: ResourceWorkspace(),
		InstanceState: map[string]string{
			"account_id":     "abc",
			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"is_no_public_ip_enabled":                  "true",
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
			"workspace_id":                             "1234",
		},
		State: map[string]interface{}{
			"account_id": "THIS_IS_CHANGING",

			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"is_no_public_ip_enabled":                  true,
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
			"workspace_id":                             1234,
		},
		Update: true,
		ID:     "abc/1234",
	}.ExpectError(t, "it is not allowed to change account_id on a running workspace")
}

func TestResourceWorkspaceUpdateLegacyConfig(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "PATCH",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
				ExpectedRequest: map[string]interface{}{
					"credentials_id": "bcd",
					"network_id":     "fgh",
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/accounts/abc/workspaces/1234",
				Response: Workspace{
					WorkspaceStatus:                     WorkspaceStatusRunning,
					IsNoPublicIPEnabled:                 true,
					WorkspaceName:                       "labdata",
					DeploymentName:                      "900150983cd24fb0",
					AwsRegion:                           "us-east-1",
					CredentialsID:                       "bcd",
					StorageConfigurationID:              "ghi",
					NetworkID:                           "fgh",
					ManagedServicesCustomerManagedKeyID: "def",
					AccountID:                           "abc",
					WorkspaceID:                         1234,
				},
			},
		},
		Resource: ResourceWorkspace(),
		State: map[string]interface{}{
			"account_id":               "abc",
			"aws_region":               "us-east-1",
			"credentials_id":           "bcd",
			"customer_managed_key_id":  "def",
			"deployment_name":          "900150983cd24fb0",
			"workspace_name":           "labdata",
			"network_id":               "fgh",
			"storage_configuration_id": "ghi",
			"workspace_id":             1234,
		},
		Update: true,
		ID:     "abc/1234",
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc/1234", d.Id(), "Id should be the same as in reading")
}

func TestResourceWorkspaceUpdate_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "PATCH",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceWorkspace(),
		State: map[string]interface{}{
			"account_id":     "abc",
			"aws_region":     "us-east-1",
			"credentials_id": "bcd",
			"managed_services_customer_managed_key_id": "def",
			"storage_customer_managed_key_id":          "def",
			"deployment_name":                          "900150983cd24fb0",
			"workspace_name":                           "labdata",
			"network_id":                               "fgh",
			"storage_configuration_id":                 "ghi",
			"workspace_id":                             1234,
		},
		Update: true,
		ID:     "abc/1234",
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "abc/1234", d.Id())
}

func TestResourceWorkspaceDelete(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "DELETE",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
				Response: Workspace{
					WorkspaceName:          "labdata",
					WorkspaceStatus:        WorkspaceStatusCanceled,
					WorkspaceStatusMessage: "Things are being removed",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
				Response: common.APIErrorBody{
					ErrorCode: "NOT_FOUND",
					Message:   "Cannot find anything",
				},
				Status: 404,
			},
		},
		Resource: ResourceWorkspace(),
		Delete:   true,
		ID:       "abc/1234",
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc/1234", d.Id())
}

func TestResourceWorkspaceDelete_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "DELETE",
				Resource: "/api/2.0/accounts/abc/workspaces/1234",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceWorkspace(),
		Delete:   true,
		ID:       "abc/1234",
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "abc/1234", d.Id())
}

func TestWaitForRunning(t *testing.T) {
	client, server, err := qa.HttpFixtureClient(t, []qa.HTTPFixture{
		{
			Method:   "POST",
			Resource: "/api/2.0/accounts/abc/workspaces",
			ExpectedRequest: Workspace{
				AccountID:                           "abc",
				IsNoPublicIPEnabled:                 true,
				WorkspaceName:                       "labdata",
				DeploymentName:                      "900150983cd24fb0",
				AwsRegion:                           "us-east-1",
				CredentialsID:                       "bcd",
				StorageConfigurationID:              "ghi",
				NetworkID:                           "fgh",
				ManagedServicesCustomerManagedKeyID: "def",
				StorageCustomerManagedKeyID:         "def",
			},
			Response: Workspace{
				WorkspaceID:    1234,
				AccountID:      "abc",
				DeploymentName: "900150983cd24fb0",
			},
		},
		{
			Method:   "GET",
			Resource: "/api/2.0/accounts/abc/workspaces/1234",
			Response: Workspace{
				WorkspaceID:                         1234,
				WorkspaceStatus:                     WorkspaceStatusProvisioning,
				WorkspaceName:                       "labdata",
				DeploymentName:                      "900150983cd24fb0",
				AwsRegion:                           "us-east-1",
				CredentialsID:                       "bcd",
				StorageConfigurationID:              "ghi",
				NetworkID:                           "fgh",
				ManagedServicesCustomerManagedKeyID: "def",
				StorageCustomerManagedKeyID:         "def",
				AccountID:                           "abc",
			},
		},
		{
			Method:   "GET",
			Resource: "/api/2.0/accounts/abc/workspaces/1234",
			Response: Workspace{
				WorkspaceID:                         1234,
				WorkspaceStatus:                     WorkspaceStatusRunning,
				WorkspaceName:                       "labdata",
				DeploymentName:                      "900150983cd24fb0",
				AwsRegion:                           "us-east-1",
				CredentialsID:                       "bcd",
				StorageConfigurationID:              "ghi",
				NetworkID:                           "fgh",
				ManagedServicesCustomerManagedKeyID: "def",
				StorageCustomerManagedKeyID:         "def",
				AccountID:                           "abc",
			},
		},
	})
	require.NoError(t, err)
	defer server.Close()

	err = NewWorkspacesAPI(context.Background(), client).Create(&Workspace{
		AccountID:                           "abc",
		IsNoPublicIPEnabled:                 true,
		WorkspaceName:                       "labdata",
		DeploymentName:                      "900150983cd24fb0",
		AwsRegion:                           "us-east-1",
		CredentialsID:                       "bcd",
		StorageConfigurationID:              "ghi",
		NetworkID:                           "fgh",
		ManagedServicesCustomerManagedKeyID: "def",
		StorageCustomerManagedKeyID:         "def",
	}, DefaultProvisionTimeout)
	require.NoError(t, err)
}

func TestCreateFailsAndCleansUp(t *testing.T) {
	client, server, err := qa.HttpFixtureClient(t, []qa.HTTPFixture{
		{
			Method:   "POST",
			Resource: "/api/2.0/accounts/abc/workspaces",
			ExpectedRequest: Workspace{
				AccountID:                           "abc",
				IsNoPublicIPEnabled:                 true,
				WorkspaceName:                       "labdata",
				DeploymentName:                      "900150983cd24fb0",
				AwsRegion:                           "us-east-1",
				CredentialsID:                       "bcd",
				StorageConfigurationID:              "ghi",
				NetworkID:                           "fgh",
				ManagedServicesCustomerManagedKeyID: "def",
				StorageCustomerManagedKeyID:         "def",
			},
			Response: Workspace{
				WorkspaceID:    1234,
				AccountID:      "abc",
				DeploymentName: "900150983cd24fb0",
			},
		},
		{
			Method:   "GET",
			Resource: "/api/2.0/accounts/abc/workspaces/1234",
			Response: Workspace{
				WorkspaceID:            1234,
				WorkspaceStatus:        WorkspaceStatusFailed,
				WorkspaceStatusMessage: "Always fails",
				WorkspaceName:          "labdata",
				DeploymentName:         "900150983cd24fb0",
				AwsRegion:              "us-east-1",
				NetworkID:              "fgh",
				AccountID:              "abc",
			},
		},
		{
			Method:   "GET",
			Resource: "/api/2.0/accounts/abc/networks/fgh",
			Response: Network{
				ErrorMessages: []NetworkHealth{
					{"FAIL", "Message"},
				},
			},
		},
		{
			Method:   "DELETE",
			Resource: "/api/2.0/accounts/abc/workspaces/1234",
		},
		{
			Method:   "GET",
			Resource: "/api/2.0/accounts/abc/workspaces/1234",
			Status:   404,
		},
	})
	require.NoError(t, err)
	defer server.Close()

	err = NewWorkspacesAPI(context.Background(), client).Create(&Workspace{
		AccountID:                           "abc",
		IsNoPublicIPEnabled:                 true,
		WorkspaceName:                       "labdata",
		DeploymentName:                      "900150983cd24fb0",
		AwsRegion:                           "us-east-1",
		CredentialsID:                       "bcd",
		StorageConfigurationID:              "ghi",
		NetworkID:                           "fgh",
		ManagedServicesCustomerManagedKeyID: "def",
		StorageCustomerManagedKeyID:         "def",
	}, DefaultProvisionTimeout)
	require.EqualError(t, err, "Workspace failed to create: Always fails, network error message: error: FAIL;error_msg: Message;")
}

func TestListWorkspaces(t *testing.T) {
	client, server, err := qa.HttpFixtureClient(t, []qa.HTTPFixture{
		{
			Method:   "GET",
			Resource: "/api/2.0/accounts/abc/workspaces",
			Response: []Workspace{},
		},
	})
	require.NoError(t, err)
	defer server.Close()

	l, err := NewWorkspacesAPI(context.Background(), client).List("abc")
	require.NoError(t, err)
	assert.Len(t, l, 0)
}

func TestDial(t *testing.T) {
	err := dial("127.0.0.1:32456", "localhost", 50*time.Millisecond)
	assert.NotNil(t, err)
	assert.Equal(t, err.Err.Error(), "dial tcp 127.0.0.1:32456: connect: connection refused")

	s := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer s.Close()
	err = dial(strings.ReplaceAll(s.URL, "http://", ""), s.URL, 500*time.Millisecond)
	assert.Nil(t, err)
}
