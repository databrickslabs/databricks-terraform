




# Examples

## Creating mount to ADLS Gen2 with AAD passthrough

(See [documentation](https://docs.microsoft.com/en-us/azure/databricks/security/credential-passthrough/adls-passthrough#--mount-azure-data-lake-storage-to-dbfs-using-credential-passthrough) for more details).

```hcl
provider "azurerm" {
  features {}
}

variable "resource_group" {
  type = string
  description = "Resource group for Databricks Workspace"
}

variable "workspace_name" {
  type = string
  description = "Name of the Databricks Workspace"
}

data "azurerm_databricks_workspace" "this" {
  name                = var.workspace_name
  resource_group_name = var.resource_group
}

# it works only with AAD token!
provider "databricks" {
  azure_workspace_resource_id = data.azurerm_databricks_workspace.this.id
}


data "databricks_node_type" "smallest" {
  local_disk = true
}

data "databricks_spark_version" "latest" {
}

resource "databricks_cluster" "shared_passthrough" {
  cluster_name            = "Shared Passthrough for mount"
  spark_version           = data.databricks_spark_version.latest.id
  node_type_id            = data.databricks_node_type.smallest.id
  autotermination_minutes = 10
  num_workers = 1
  
  spark_conf = {
    "spark.databricks.cluster.profile":"serverless",
    "spark.databricks.repl.allowedLanguages":"python,sql",
    "spark.databricks.passthrough.enabled": "true",
    "spark.databricks.pyspark.enableProcessIsolation": "true"
  }
  
  custom_tags = {
    "ResourceClass": "Serverless"
  }
}

variable "storage_acc" {
  type = string
  description = "Name of the ADLS Gen2 storage container"
}

variable "container" {
  type = string
  description = "Name of container inside storage account"
}

resource "databricks_mount" "passthrough" {
  mount_name = "passthrough-test"
  cluster_id = databricks_cluster.shared_passthrough.id
  
  source = "abfss://${var.container}@${var.storage_acc}.dfs.core.windows.net"
  extra_configs = {
    "fs.azure.account.auth.type": "CustomAccessToken",
    "fs.azure.account.custom.token.provider.class": "{{sparkconf/spark.databricks.passthrough.adls.gen2.tokenProviderClassName}}",
  }
}
```
