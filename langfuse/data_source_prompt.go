package langfuse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrompt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePromptRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prompt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"commit_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_out": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcePromptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	name := d.Get("name").(string)

	query := url.Values{}
	if v, ok := d.GetOk("version"); ok {
		query.Set("version", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("label"); ok {
		query.Set("label", v.(string))
	}

	path := fmt.Sprintf("/api/public/v2/prompts/%s", name)
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	var resp map[string]interface{}
	if err := client.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return diag.FromErr(err)
	}

	if t, ok := resp["type"].(string); ok {
		if err := d.Set("type", t); err != nil {
			return diag.FromErr(err)
		}
	}

	if p, ok := resp["prompt"]; ok {
		b, err := json.Marshal(p)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("prompt", string(b)); err != nil {
			return diag.FromErr(err)
		}
	}

	if cfg, ok := resp["config"].(map[string]interface{}); ok {
		if err := d.Set("config", cfg); err != nil {
			return diag.FromErr(err)
		}
	}

	if labels, ok := resp["labels"].([]interface{}); ok {
		if err := d.Set("labels", labels); err != nil {
			return diag.FromErr(err)
		}
	}

	if tags, ok := resp["tags"].([]interface{}); ok {
		if err := d.Set("tags", tags); err != nil {
			return diag.FromErr(err)
		}
	}

	if msg, ok := resp["commitMessage"].(string); ok {
		if err := d.Set("commit_message", msg); err != nil {
			return diag.FromErr(err)
		}
	}

	version, ok := resp["version"].(float64)
	if ok {
		ver := int(version)
		d.SetId(fmt.Sprintf("%s:%d", name, ver))
		if err := d.Set("version_out", ver); err != nil {
			return diag.FromErr(err)
		}
	} else {
		d.SetId(name)
	}

	return nil
}
