package databricks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/databrickslabs/databricks-terraform/client/service"
)

// List of management information
const (
	ADBResourceID string = "2ff814a6-3304-4ab8-85cb-cd0e6f879c1d"
)

// AzureAuth is a struct that contains information about the azure sp authentication
type AzureAuth struct {
	TokenPayload           *TokenPayload
	ManagementToken        string
	AdbWorkspaceResourceID string
	AdbAccessToken         string
	AdbPlatformToken       string
	AdbWorkspaceUrl        string
}

// TokenPayload contains all the auth information for azure sp authentication
type TokenPayload struct {
	ManagedResourceGroup string
	AzureRegion          string
	WorkspaceName        string
	ResourceGroup        string
	SubscriptionID       string
	ClientSecret         string
	ClientID             string
	TenantID             string
}

// WsProps contains information about the workspace properties
type WsProps struct {
	ManagedResourceGroupID string `json:"managedResourceGroupId"`
}

// WorkspaceRequest contains the request information for getting workspace information
type WorkspaceRequest struct {
	Properties *WsProps `json:"properties"`
	Name       string   `json:"name"`
	Location   string   `json:"location"`
}

func (a *AzureAuth) getManagementToken() error {
	log.Println("[DEBUG] Creating Azure Databricks management OAuth token.")
	mgmtTokenOAuthCfg, err := adal.NewOAuthConfigWithAPIVersion(azure.PublicCloud.ActiveDirectoryEndpoint,
		a.TokenPayload.TenantID,
		nil)
	if err != nil {
		return err
	}
	mgmtToken, err := adal.NewServicePrincipalToken(
		*mgmtTokenOAuthCfg,
		a.TokenPayload.ClientID,
		a.TokenPayload.ClientSecret,
		azure.PublicCloud.ServiceManagementEndpoint)
	if err != nil {
		return err
	}
	err = mgmtToken.Refresh()
	if err != nil {
		return err
	}
	a.ManagementToken = mgmtToken.OAuthToken()
	return nil
}

func (a *AzureAuth) getADBPlatformToken() error {
	log.Println("[DEBUG] Creating Azure Databricks management OAuth token.")
	platformTokenOAuthCfg, err := adal.NewOAuthConfigWithAPIVersion(azure.PublicCloud.ActiveDirectoryEndpoint,
		a.TokenPayload.TenantID,
		nil)
	if err != nil {
		return err
	}
	platformToken, err := adal.NewServicePrincipalToken(
		*platformTokenOAuthCfg,
		a.TokenPayload.ClientID,
		a.TokenPayload.ClientSecret,
		ADBResourceID)
	if err != nil {
		return err
	}

	err = platformToken.Refresh()
	if err != nil {
		return err
	}
	a.AdbPlatformToken = platformToken.OAuthToken()
	return nil
}

func (a *AzureAuth) getWorkspaceAccessToken(config *service.DBApiClientConfig) error {
	log.Println("[DEBUG] Creating workspace token")
	apiLifeTimeInSeconds := int32(3600)
	comment := "Secret made via SP"
	url := a.AdbWorkspaceUrl + "/api/2.0/token/create"
	payload := struct {
		LifetimeSeconds int32  `json:"lifetime_seconds,omitempty"`
		Comment         string `json:"comment,omitempty"`
	}{
		apiLifeTimeInSeconds,
		comment,
	}
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"X-Databricks-Azure-Workspace-Resource-Id": a.AdbWorkspaceResourceID,
		"X-Databricks-Azure-SP-Management-Token":   a.ManagementToken,
		"cache-control":                            "no-cache",
		"Authorization":                            "Bearer " + a.AdbPlatformToken,
	}

	var responseMap map[string]interface{}
	resp, err := service.PerformQuery(config, http.MethodPost, url, "2.0", headers, true, true, payload, nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &responseMap)
	if err != nil {
		return err
	}
	a.AdbAccessToken = responseMap["token_value"].(string)
	return nil
}

// Main function call that gets made and it follows 4 steps at the moment:
// 1. Get Management OAuth Token using management endpoint
// 2. Get Workspace ID
// 3. Get Azure Databricks Platform OAuth Token using Databricks resource id
// 4. Get Azure Databricks Workspace Personal Access Token for the SP (60 min duration)
func (a *AzureAuth) initWorkspaceAndGetClient(config *service.DBApiClientConfig) error {
	//var dbClient service.DBApiClient

	// Get management token
	err := a.getManagementToken()
	if err != nil {
		return err
	}

	a.AdbWorkspaceResourceID = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Databricks/workspaces/%s",
		a.TokenPayload.SubscriptionID,
		a.TokenPayload.ResourceGroup,
		a.TokenPayload.WorkspaceName)

	// Get platform token
	err = a.getADBPlatformToken()
	if err != nil {
		return err
	}

	// Get workspace personal access token
	err = a.getWorkspaceAccessToken(config)
	if err != nil {
		return err
	}

	config.Host = a.AdbWorkspaceUrl
	config.Token = a.AdbAccessToken

	return nil
}
