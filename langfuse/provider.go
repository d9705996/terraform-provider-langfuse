package langfuse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderConfig struct {
	Host     string
	Username string
	Password string
}

type apiClient struct {
	baseURL    *url.URL
	httpClient *http.Client
	username   string
	password   string
}

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://cloud.langfuse.com",
				Description: "Base URL for Langfuse API",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LANGFUSE_USERNAME", nil),
				Description: "Langfuse public key",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("LANGFUSE_PASSWORD", nil),
				Description: "Langfuse secret key",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"langfuse_project":         resourceProject(),
			"langfuse_project_api_key": resourceProjectAPIKey(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"langfuse_project":          dataSourceProject(),
			"langfuse_project_api_keys": dataSourceProjectApiKeys(),
		},
	}

	p.ConfigureContextFunc = configureProvider
	return p
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	parsed, err := url.Parse(host)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid host",
			Detail:   fmt.Sprintf("failed to parse host: %v", err),
		})
		return nil, diags
	}

	client := &apiClient{
		baseURL:    parsed,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		username:   username,
		password:   password,
	}

	return client, diags
}

func (c *apiClient) doRequest(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	u := *c.baseURL
	u.Path = path

	var req *http.Request
	var err error
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req, err = http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)
		if err != nil {
			return err
		}
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			fmt.Printf("error closing response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", b)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
