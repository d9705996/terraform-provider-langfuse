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

type prompt struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func resourcePrompt() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePromptCreate,
		ReadContext:   resourcePromptRead,
		UpdateContext: resourcePromptUpdate,
		DeleteContext: resourcePromptDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePromptImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePromptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pid := d.Get("project_id").(string)
	body := map[string]interface{}{
		"content": d.Get("content").(string),
	}

	var resp prompt
	path := fmt.Sprintf("/api/public/projects/%s/prompts", pid)
	if err := client.doRequest(ctx, http.MethodPost, path, body, &resp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	if err := d.Set("content", resp.Content); err != nil {
		return diag.FromErr(err)
	}
	return resourcePromptRead(ctx, d, meta)
}

func resourcePromptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pid := d.Get("project_id").(string)

	var resp prompt
	path := fmt.Sprintf("/api/public/projects/%s/prompts/%s", pid, d.Id())
	if err := client.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		if strings.Contains(err.Error(), "404") {
			tflog.Trace(ctx, "prompt not found, removing from state")
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("content", resp.Content); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePromptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pid := d.Get("project_id").(string)
	body := map[string]interface{}{
		"content": d.Get("content").(string),
	}

	var resp prompt
	path := fmt.Sprintf("/api/public/projects/%s/prompts/%s", pid, d.Id())
	if err := client.doRequest(ctx, http.MethodPut, path, body, &resp); err != nil {
		return diag.FromErr(err)
	}

	return resourcePromptRead(ctx, d, meta)
}

func resourcePromptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pid := d.Get("project_id").(string)
	path := fmt.Sprintf("/api/public/projects/%s/prompts/%s", pid, d.Id())
	if err := client.doRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "deleted prompt")
	d.SetId("")
	return nil
}

func resourcePromptImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected project_id/prompt_id", d.Id())
	}
	if err := d.Set("project_id", parts[0]); err != nil {
		return nil, err
	}
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
