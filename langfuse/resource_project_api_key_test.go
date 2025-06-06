package langfuse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectAPIKeyResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			resp := map[string]interface{}{
				"id":               "1",
				"createdAt":        "2024-01-01T00:00:00Z",
				"publicKey":        "pk",
				"secretKey":        "sk",
				"displaySecretKey": "disp-sk",
				"note":             "ci",
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("encode response: %v", err)
			}
		case http.MethodGet:
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
		case http.MethodDelete:
			if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
				t.Fatalf("encode response: %v", err)
			}
		}
	}))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectAPIKey(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("langfuse_project_api_key.test", "public_key", "pk"),
					resource.TestCheckResourceAttr("langfuse_project_api_key.test", "note", "ci"),
				),
			},
		},
	})
}

func testAccProjectAPIKey(url string) string {
	return fmt.Sprintf(`
provider "langfuse" {
  host = "%s"
}

resource "langfuse_project_api_key" "test" {
  project_id = "123"
  note       = "ci"
}
`, url)
}
