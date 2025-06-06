package langfuse

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type project struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Metadata      map[string]interface{} `json:"metadata"`
	RetentionDays int                    `json:"retentionDays"`
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	body := map[string]interface{}{
		"name":      d.Get("name").(string),
		"metadata":  d.Get("metadata"),
		"retention": d.Get("retention").(int),
	}

	var resp project
	if err := client.doRequest(ctx, http.MethodPost, "/api/public/projects", body, &resp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	var resp project
	path := fmt.Sprintf("/api/public/projects/%s", d.Id())
	if err := client.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", resp.Name)
	d.Set("metadata", resp.Metadata)
	d.Set("retention", resp.RetentionDays)

	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	body := map[string]interface{}{
		"name":      d.Get("name").(string),
		"metadata":  d.Get("metadata"),
		"retention": d.Get("retention").(int),
	}

	var resp project
	path := fmt.Sprintf("/api/public/projects/%s", d.Id())
	if err := client.doRequest(ctx, http.MethodPut, path, body, &resp); err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	path := fmt.Sprintf("/api/public/projects/%s", d.Id())
	if err := client.doRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "deleted project")
	d.SetId("")
	return nil
}
