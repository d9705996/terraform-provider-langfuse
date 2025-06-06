package langfuse

import (
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
