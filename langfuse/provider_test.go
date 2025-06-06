package langfuse

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func testProvider() *schema.Provider {
	return Provider()
}

var providerFactories = map[string]func() (*schema.Provider, error){
	"langfuse": func() (*schema.Provider, error) { return testProvider(), nil },
}

func testAccPreCheck(t *testing.T) {}

func TestProviderHostEnvVar(t *testing.T) {
	const h = "http://example.com"
	if err := os.Setenv("LANGFUSE_HOST", h); err != nil {
		t.Fatalf("set env: %v", err)
	}
	defer func() { _ = os.Unsetenv("LANGFUSE_HOST") }()

	p := testProvider()
	d := schema.TestResourceDataRaw(t, p.Schema, map[string]interface{}{})

	raw, diags := p.ConfigureContextFunc(context.Background(), d)
	if len(diags) != 0 {
		t.Fatalf("unexpected diags: %v", diags)
	}

	c := raw.(*apiClient)
	if c.baseURL.String() != h {
		t.Fatalf("expected host %q, got %q", h, c.baseURL.String())
	}
}
