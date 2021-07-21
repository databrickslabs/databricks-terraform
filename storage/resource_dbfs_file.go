package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/databrickslabs/terraform-provider-databricks/common"

	"github.com/databrickslabs/terraform-provider-databricks/workspace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDBFSFile manages files on DBFS
func ResourceDBFSFile() *schema.Resource {
	return common.Resource{
		SchemaVersion: 1,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, c interface{}) error {
			a, b := d.GetChange("modification_time")
			log.Printf("[DEBUG] ---------> modification time diff %v vs %v", a, b)
			if d.HasChange("modification_time") {
				log.Printf("[DEBUG] ---------> detected change in modification_time")
				return d.SetNewComputed("modification_time")
			}
			return nil
		},
		Schema: workspace.FileContentSchema(map[string]*schema.Schema{
			"file_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"modification_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"dbfs_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
		Create: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			path := d.Get("path").(string)
			content, err := workspace.ReadContent(d)
			if err != nil {
				return err
			}
			if err = NewDbfsAPI(ctx, c).Create(path, content, true); err != nil {
				return err
			}
			d.SetId(path)
			return nil
		},
		Read: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			dbfsAPI := NewDbfsAPI(ctx, c)
			fileInfo, err := dbfsAPI.Status(d.Id())
			if err != nil {
				return err
			}
			d.Set("path", fileInfo.Path)
			d.Set("dbfs_path", fmt.Sprint("dbfs:", fileInfo.Path))
			d.Set("file_size", fileInfo.FileSize)
			d.Set("modification_time", fileInfo.ModificationTime)
			return nil
		},
		Delete: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			return NewDbfsAPI(ctx, c).Delete(d.Id(), false)
		},
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    DbfsFileV0(),
				Upgrade: workspace.MigrateV0,
			},
		},
	}.ToResource()
}
