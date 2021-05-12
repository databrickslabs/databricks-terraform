package provider

import (
	"context"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/databrickslabs/terraform-provider-databricks/access"
	"github.com/databrickslabs/terraform-provider-databricks/common"
	"github.com/databrickslabs/terraform-provider-databricks/compute"
	"github.com/databrickslabs/terraform-provider-databricks/identity"
	"github.com/databrickslabs/terraform-provider-databricks/mws"
	"github.com/databrickslabs/terraform-provider-databricks/sqlanalytics"
	"github.com/databrickslabs/terraform-provider-databricks/storage"
	"github.com/databrickslabs/terraform-provider-databricks/workspace"
)

// DatabricksProvider returns the entire terraform provider object
func DatabricksProvider() *schema.Provider {
	p := &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"databricks_aws_crossaccount_policy": access.DataAwsCrossAccountPolicy(),
			"databricks_aws_assume_role_policy":  access.DataAwsAssumeRolePolicy(),
			"databricks_aws_bucket_policy":       access.DataAwsBucketPolicy(),
			"databricks_current_user":            identity.DataSourceCurrentUser(),
			"databricks_dbfs_file":               storage.DataSourceDBFSFile(),
			"databricks_dbfs_file_paths":         storage.DataSourceDBFSFilePaths(),
			"databricks_group":                   identity.DataSourceGroup(),
			"databricks_node_type":               compute.DataSourceNodeType(),
			"databricks_notebook":                workspace.DataSourceNotebook(),
			"databricks_notebook_paths":          workspace.DataSourceNotebookPaths(),
			"databricks_spark_version":           compute.DataSourceSparkVersion(),
			"databricks_user":                    identity.DataSourceUser(),
			"databricks_zones":                   compute.DataSourceClusterZones(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"databricks_secret":          access.ResourceSecret(),
			"databricks_secret_scope":    access.ResourceSecretScope(),
			"databricks_secret_acl":      access.ResourceSecretACL(),
			"databricks_permissions":     access.ResourcePermissions(),
			"databricks_sql_permissions": access.ResourceSqlPermissions(),
			"databricks_ip_access_list":  access.ResourceIPAccessList(),

			"databricks_cluster":        compute.ResourceCluster(),
			"databricks_cluster_policy": compute.ResourceClusterPolicy(),
			"databricks_instance_pool":  compute.ResourceInstancePool(),
			"databricks_job":            compute.ResourceJob(),
			"databricks_pipeline":       compute.ResourcePipeline(),

			"databricks_group":                  identity.ResourceGroup(),
			"databricks_group_instance_profile": identity.ResourceGroupInstanceProfile(),
			"databricks_user_instance_profile":  identity.ResourceUserInstanceProfile(),
			"databricks_instance_profile":       identity.ResourceInstanceProfile(),
			"databricks_group_member":           identity.ResourceGroupMember(),
			"databricks_token":                  identity.ResourceToken(),
			"databricks_user":                   identity.ResourceUser(),
			"databricks_service_principal":      identity.ResourceServicePrincipal(),

			"databricks_mws_customer_managed_keys":   mws.ResourceCustomerManagedKey(),
			"databricks_mws_credentials":             mws.ResourceCredentials(),
			"databricks_mws_log_delivery":            mws.ResourceLogDelivery(),
			"databricks_mws_networks":                mws.ResourceNetwork(),
			"databricks_mws_private_access_settings": mws.ResourcePrivateAccessSettings(),
			"databricks_mws_storage_configurations":  mws.ResourceStorageConfiguration(),
			"databricks_mws_vpc_endpoint":            mws.ResourceVPCEndpoint(),
			"databricks_mws_workspaces":              mws.ResourceWorkspace(),

			"databricks_aws_s3_mount":          storage.ResourceAWSS3Mount(),
			"databricks_azure_adls_gen1_mount": storage.ResourceAzureAdlsGen1Mount(),
			"databricks_azure_adls_gen2_mount": storage.ResourceAzureAdlsGen2Mount(),
			"databricks_azure_blob_mount":      storage.ResourceAzureBlobMount(),
			"databricks_dbfs_file":             storage.ResourceDBFSFile(),
			"databricks_mount":                 storage.ResourceDatabricksMount(),

			"databricks_sql_dashboard":     sqlanalytics.ResourceDashboard(),
			"databricks_sql_endpoint":      sqlanalytics.ResourceSQLEndpoint(),
			"databricks_sql_query":         sqlanalytics.ResourceQuery(),
			"databricks_sql_visualization": sqlanalytics.ResourceVisualization(),
			"databricks_sql_widget":        sqlanalytics.ResourceWidget(),

			"databricks_directory":          workspace.ResourceDirectory(),
			"databricks_global_init_script": workspace.ResourceGlobalInitScript(),
			"databricks_notebook":           workspace.ResourceNotebook(),
			"databricks_workspace_conf":     workspace.ResourceWorkspaceConf(),
		},
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_HOST", nil),
				ConflictsWith: []string{
					"config_file",
				},
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_TOKEN", nil),
				ConflictsWith: []string{
					"username",
					"password",
					"config_file",
					"profile",
				},
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_USERNAME", nil),
				ConflictsWith: []string{
					"token",
				},
			},
			"password": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_PASSWORD", nil),
				ConflictsWith: []string{
					"token",
				},
			},
			"config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_CONFIG_FILE", nil),
				Description: "Location of the Databricks CLI credentials file, that is created\n" +
					"by `databricks configure --token` command. By default, it is located\n" +
					"in ~/.databrickscfg. Check  https://docs.databricks.com/dev-tools/cli/index.html#set-up-authentication for docs. Config\n" +
					"file credentials will only be used when host/token are not provided.",
				ConflictsWith: []string{
					"token",
					"host",
				},
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_CONFIG_PROFILE", nil),
				Description: "Connection profile specified within ~/.databrickscfg. Please check\n" +
					"https://docs.databricks.com/dev-tools/cli/index.html#connection-profiles for documentation.",
				ConflictsWith: []string{
					"token",
				},
			},
			"azure_workspace_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"DATABRICKS_AZURE_WORKSPACE_RESOURCE_ID",
					"AZURE_DATABRICKS_WORKSPACE_RESOURCE_ID"}, nil),
				ConflictsWith: []string{
					"azure_workspace_name",
				},
			},
			"azure_workspace_name": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("DATABRICKS_AZURE_WORKSPACE_NAME", nil),
				ConflictsWith: []string{"azure_workspace_resource_id"},
			},
			"azure_resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_AZURE_RESOURCE_GROUP", nil),
			},
			"azure_subscription_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"DATABRICKS_AZURE_SUBSCRIPTION_ID",
					"ARM_SUBSCRIPTION_ID"}, nil),
			},
			"azure_client_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"DATABRICKS_AZURE_CLIENT_ID",
					"ARM_CLIENT_ID"}, nil),
			},
			"azure_client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"DATABRICKS_AZURE_CLIENT_SECRET",
					"ARM_CLIENT_SECRET"}, nil),
			},
			"azure_tenant_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"DATABRICKS_AZURE_TENANT_ID",
					"ARM_TENANT_ID"}, nil),
			},
			"azure_pat_token_duration_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Currently secret scopes are not accessible via AAD tokens so we will need to create a PAT token",
				Default:     "3600",
			},
			"azure_use_pat_for_cli": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Create ephemeral PAT tokens also for AZ CLI authenticated requests",
			},
			"azure_use_pat_for_spn": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Create ephemeral PAT tokens instead of AAD tokens for SPN",
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_AZURE_USE_PAT_FOR_SPN", false),
			},
			"azure_environment": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
			},
			"skip_verify": {
				Type:        schema.TypeBool,
				Description: "Skip SSL certificate verification for HTTP calls. Use at your own risk.",
				Optional:    true,
				Default:     false,
			},
			"debug_truncate_bytes": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Truncate JSON fields in JSON above this limit. Default is 96. Visible only when TF_LOG=DEBUG is set",
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_DEBUG_TRUNCATE_BYTES", common.DefaultTruncateBytes),
			},
			"debug_headers": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Debug HTTP headers of requests made by the provider. Default is false. Visible only when TF_LOG=DEBUG is set",
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_DEBUG_HEADERS", false),
			},
			"rate_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum number of requests per second made to Databricks REST API by Terraform.",
				DefaultFunc: schema.EnvDefaultFunc("DATABRICKS_RATE_LIMIT", common.DefaultRateLimitPerSecond),
			},
		},
	}
	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		ctx = context.WithValue(ctx, common.Provider, p)
		return configureDatabricksClient(ctx, d)
	}
	common.AddContextToAllResources(p, "databricks")
	return p
}

func configureDatabricksClient(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	prov := ctx.Value(common.Provider).(*schema.Provider)
	pc := common.DatabricksClient{
		Provider: prov,
	}
	attrsUsed := []string{}
	for attr, schm := range prov.Schema {
		if schm.Type != schema.TypeString {
			continue
		}
		if attr == "azure_pat_token_duration_seconds" {
			continue
		}
		if value, ok := d.GetOk(attr); ok {
			if value.(string) == "" {
				continue
			}
			if attr == "azure_environment" && value.(string) == "public" {
				continue
			}
			attrsUsed = append(attrsUsed, attr)
		}
	}
	sort.Strings(attrsUsed)
	log.Printf("[INFO] Explicit and implicit attributes: %s", strings.Join(attrsUsed, ", "))
	authsUsed := map[string]bool{}
	if host, ok := d.GetOk("host"); ok {
		pc.Host = host.(string)
	}
	if token, ok := d.GetOk("token"); ok {
		authsUsed["token"] = true
		pc.Token = token.(string)
	}
	if v, ok := d.GetOk("username"); ok {
		authsUsed["password"] = true
		pc.Username = v.(string)
	}
	if v, ok := d.GetOk("password"); ok {
		authsUsed["password"] = true
		pc.Password = v.(string)
	}
	if v, ok := d.GetOk("profile"); ok {
		authsUsed["config profile"] = true
		pc.Profile = v.(string)
	}
	if v, ok := d.GetOk("config_file"); ok {
		authsUsed["config profile"] = true
		pc.ConfigFile = v.(string)
	}
	if v, ok := d.GetOk("azure_workspace_resource_id"); ok {
		authsUsed["azure"] = true
		pc.AzureAuth.ResourceID = v.(string)
	}
	if v, ok := d.GetOk("azure_workspace_name"); ok {
		authsUsed["azure"] = true
		pc.AzureAuth.WorkspaceName = v.(string)
	}
	if v, ok := d.GetOk("azure_resource_group"); ok {
		authsUsed["azure"] = true
		pc.AzureAuth.ResourceGroup = v.(string)
	}
	if v, ok := d.GetOk("azure_subscription_id"); ok {
		authsUsed["azure"] = true
		pc.AzureAuth.SubscriptionID = v.(string)
	}
	if v, ok := d.GetOk("azure_client_secret"); ok {
		authsUsed["azure"] = true
		pc.AzureAuth.ClientSecret = v.(string)
	}
	if v, ok := d.GetOk("azure_client_id"); ok {
		authsUsed["azure"] = true
		pc.AzureAuth.ClientID = v.(string)
	}
	if v, ok := d.GetOk("azure_tenant_id"); ok {
		authsUsed["azure"] = true
		pc.AzureAuth.TenantID = v.(string)
	}
	if v, ok := d.GetOk("azure_pat_token_duration_seconds"); ok {
		pc.AzureAuth.PATTokenDurationSeconds = v.(string)
	}
	if v, ok := d.GetOk("skip_verify"); ok {
		pc.InsecureSkipVerify = v.(bool)
	}
	if v, ok := d.GetOk("debug_truncate_bytes"); ok {
		pc.DebugTruncateBytes = v.(int)
	}
	if v, ok := d.GetOk("rate_limit"); ok {
		pc.RateLimitPerSecond = v.(int)
	}
	if v, ok := d.GetOk("debug_headers"); ok {
		pc.DebugHeaders = v.(bool)
	}
	if v, ok := d.GetOk("azure_use_pat_for_cli"); ok {
		pc.AzureAuth.UsePATForCLI = v.(bool)
	}
	if v, ok := d.GetOk("azure_use_pat_for_spn"); ok {
		pc.AzureAuth.UsePATForSPN = v.(bool)
	}
	if v, ok := d.GetOk("azure_environment"); ok {
		pc.AzureAuth.Environment = v.(string)
	}
	authorizationMethodsUsed := []string{}
	for name, used := range authsUsed {
		if used {
			authorizationMethodsUsed = append(authorizationMethodsUsed, name)
		}
	}
	if len(authorizationMethodsUsed) > 1 {
		sort.Strings(authorizationMethodsUsed)
		return nil, diag.Errorf("More than one authorization method configured: %s",
			strings.Join(authorizationMethodsUsed, " and "))
	}
	if err := pc.Configure(); err != nil {
		return nil, diag.FromErr(err)
	}
	pc.WithCommandExecutor(func(ctx context.Context, client *common.DatabricksClient) common.CommandExecutor {
		return compute.NewCommandsAPI(ctx, client)
	})
	return &pc, nil
}
