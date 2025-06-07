package langfuse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectAPIKeysDataSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		resp := map[string]interface{}{
			"apiKeys": []map[string]interface{}{
				{
					"id":               "1",
					"createdAt":        "2024-01-01T00:00:00Z",
					"publicKey":        "pk",
					"displaySecretKey": "disp-sk",
					"note":             "ci",
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjectAPIKeys(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.langfuse_project_api_keys.test", "api_keys.0.id", "1"),
					resource.TestCheckResourceAttr("data.langfuse_project_api_keys.test", "api_keys.0.public_key", "pk"),
				),
			},
		},
	})
}

func testAccDataSourceProjectAPIKeys(url string) string {
	return fmt.Sprintf(`
provider "langfuse" {
  host = "%s"
}

data "langfuse_project_api_keys" "test" {
  project_id = "123"
}
`, url)
}
