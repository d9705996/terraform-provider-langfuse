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

type apiKeyResponse struct {
	ID               string  `json:"id"`
	CreatedAt        string  `json:"createdAt"`
	PublicKey        string  `json:"publicKey"`
	SecretKey        string  `json:"secretKey"`
	DisplaySecretKey string  `json:"displaySecretKey"`
	Note             *string `json:"note"`
}

type apiKeySummary struct {
	ID               string  `json:"id"`
	CreatedAt        string  `json:"createdAt"`
	ExpiresAt        *string `json:"expiresAt"`
	LastUsedAt       *string `json:"lastUsedAt"`
	Note             *string `json:"note"`
	PublicKey        string  `json:"publicKey"`
	DisplaySecretKey string  `json:"displaySecretKey"`
}

type apiKeyList struct {
	ApiKeys []apiKeySummary `json:"apiKeys"`
}

func resourceProjectAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectAPIKeyCreate,
		ReadContext:   resourceProjectAPIKeyRead,
		DeleteContext: resourceProjectAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceProjectAPIKeyImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"display_secret_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"last_used_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceProjectAPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pid := d.Get("project_id").(string)
	body := map[string]interface{}{}
	if v, ok := d.GetOk("note"); ok {
		body["note"] = v.(string)
	}

	var resp apiKeyResponse
	path := fmt.Sprintf("/api/public/projects/%s/apiKeys", pid)
	if err := client.doRequest(ctx, http.MethodPost, path, body, &resp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	if err := d.Set("public_key", resp.PublicKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secret_key", resp.SecretKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_secret_key", resp.DisplaySecretKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_at", resp.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if resp.Note != nil {
		if err := d.Set("note", *resp.Note); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceProjectAPIKeyRead(ctx, d, meta)
}

func resourceProjectAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pid := d.Get("project_id").(string)

	var list apiKeyList
	path := fmt.Sprintf("/api/public/projects/%s/apiKeys", pid)
	if err := client.doRequest(ctx, http.MethodGet, path, nil, &list); err != nil {
		return diag.FromErr(err)
	}

	var found *apiKeySummary
	for _, k := range list.ApiKeys {
		if k.ID == d.Id() {
			found = &k
			break
		}
	}
	if found == nil {
		tflog.Trace(ctx, "api key not found, removing from state")
		d.SetId("")
		return nil
	}

	if err := d.Set("public_key", found.PublicKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_secret_key", found.DisplaySecretKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_at", found.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if found.Note != nil {
		if err := d.Set("note", *found.Note); err != nil {
			return diag.FromErr(err)
		}
	}
	if found.ExpiresAt != nil {
		if err := d.Set("expires_at", *found.ExpiresAt); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("expires_at", ""); err != nil {
			return diag.FromErr(err)
		}
	}
	if found.LastUsedAt != nil {
		if err := d.Set("last_used_at", *found.LastUsedAt); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("last_used_at", ""); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceProjectAPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pid := d.Get("project_id").(string)
	path := fmt.Sprintf("/api/public/projects/%s/apiKeys/%s", pid, d.Id())
	if err := client.doRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceProjectAPIKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected project_id/api_key_id", d.Id())
	}
	if err := d.Set("project_id", parts[0]); err != nil {
		return nil, err
	}
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
