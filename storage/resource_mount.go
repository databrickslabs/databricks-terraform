package storage

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GenericMount struct {
	URI     string            `json:"source"`
	Options map[string]string `json:"extra_configs"`
}

// Source returns URI backing the mount
func (m GenericMount) Source() string {
	return m.URI
}

// Config returns mount configurations
func (m GenericMount) Config() map[string]string {
	return m.Options
}

// TODO: https://databricks.slack.com/archives/D019B5E3RBJ/p1620837190043900
// TODO: add a abfss / s3 / other blocks for mount special configs

// ResourceDatabricksMount mounts using given configuration
func ResourceDatabricksMount() *schema.Resource {
	return commonMountResource(GenericMount{}, map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"mount_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"source": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"extra_configs": {
			Type:     schema.TypeMap,
			Required: true,
			ForceNew: true,
		},
	})
}
