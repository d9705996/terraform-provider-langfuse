package langfuse

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProjectApiKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectApiKeysRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":                 {Type: schema.TypeString, Computed: true},
						"created_at":         {Type: schema.TypeString, Computed: true},
						"expires_at":         {Type: schema.TypeString, Computed: true},
						"last_used_at":       {Type: schema.TypeString, Computed: true},
						"note":               {Type: schema.TypeString, Computed: true},
						"public_key":         {Type: schema.TypeString, Computed: true},
						"display_secret_key": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceProjectApiKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pid := d.Get("project_id").(string)

	var list apiKeyList
	path := fmt.Sprintf("/api/public/projects/%s/apiKeys", pid)
	if err := client.doRequest(ctx, http.MethodGet, path, nil, &list); err != nil {
		return diag.FromErr(err)
	}

	var result []map[string]interface{}
	for _, k := range list.ApiKeys {
		m := map[string]interface{}{
			"id":                 k.ID,
			"created_at":         k.CreatedAt,
			"public_key":         k.PublicKey,
			"display_secret_key": k.DisplaySecretKey,
		}
		if k.Note != nil {
			m["note"] = *k.Note
		}
		if k.ExpiresAt != nil {
			m["expires_at"] = *k.ExpiresAt
		}
		if k.LastUsedAt != nil {
			m["last_used_at"] = *k.LastUsedAt
		}
		result = append(result, m)
	}

	d.SetId(pid)
	if err := d.Set("api_keys", result); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
