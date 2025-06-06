package langfuse

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	var resp project
	path := fmt.Sprintf("/api/public/projects/%s", d.Get("id").(string))
	if err := client.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", resp.Metadata); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("retention", resp.RetentionDays); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
