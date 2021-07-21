package storage

import (
	"github.com/databrickslabs/terraform-provider-databricks/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GenericMount struct {
	URI     string              `json:"uri,omitempty"`
	Options map[string]string   `json:"extra_configs,omitempty"`
	Abfss   *AzureADLSGen2Mount `json:"abfss,omitempty"`
	S3      *AWSIamMount        `json:"s3,omitempty"`
	Adl     *AzureADLSGen1Mount `json:"adl,omitempty"`
	Wasbs   *AzureBlobMount     `json:"wasbs,omitempty"`
}

// Source returns URI backing the mount
func (m GenericMount) Source() string {
	if m.Abfss != nil {
		return m.Abfss.Source()
	} else if m.Adl != nil {
		m.Adl.Source()
	} else if m.Wasbs != nil {
		return m.Wasbs.Source()
	} else if m.S3 != nil {
		return m.S3.Source()
	}
	return m.URI
}

// Config returns mount configurations
func (m GenericMount) Config() map[string]string {
	if m.Abfss != nil {
		return m.Abfss.Config()
	} else if m.Adl != nil {
		m.Adl.Config()
	} else if m.Wasbs != nil {
		return m.Wasbs.Config()
	} else if m.S3 != nil {
		return m.S3.Config()
	}
	return m.Options
}

func extractFieldNames(elem interface{}) []string {
	m := elem.(*schema.Resource)
	keys := make([]string, 0, len(m.Schema))
	for k := range m.Schema {
		keys = append(keys, k)
	}
	return keys
}

// ResourceDatabricksMount mounts using given configuration
func ResourceDatabricksMount() *schema.Resource {
	scm := common.StructToSchema(GenericMount{}, func(s map[string]*schema.Schema) map[string]*schema.Schema {
		s["cluster_id"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		}
		s["mount_name"] = &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		}
		s["source"] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
		s["uri"].ConflictsWith = []string{"abfss", "wasbs", "s3", "adl"}
		s["extra_configs"].ConflictsWith = []string{"abfss", "wasbs", "s3", "adl"}
		s["abfss"].ConflictsWith = []string{"uri", "extra_configs", "wasbs", "s3", "adl"}
		s["wasbs"].ConflictsWith = []string{"uri", "extra_configs", "abfss", "s3", "adl"}
		s["s3"].ConflictsWith = []string{"uri", "extra_configs", "wasbs", "abfss", "adl"}
		s["adl"].ConflictsWith = []string{"uri", "extra_configs", "wasbs", "s3", "abfss"}
		s["uri"].ForceNew = true
		s["extra_configs"].ForceNew = true
		blocks := []string{"abfss", "wasbs", "s3", "adl"}
		for _, nm := range blocks {
			s[nm].DiffSuppressFunc = common.MakeEmptyBlockSuppressFunc(nm)
			s[nm].ForceNew = true
			for _, field := range extractFieldNames(s[nm].Elem) {
				p, err := common.SchemaPath(s, nm, field)
				if err == nil {
					p.ForceNew = true
				}
			}
		}
		// TODO: We need to have a validation function that will check that source isn't empty if other blocks aren't specified

		return s
	})

	return commonMountResource(GenericMount{}, scm)
}
