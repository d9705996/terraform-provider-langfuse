package langfuse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req map[string]interface{}
			json.NewDecoder(r.Body).Decode(&req)
			resp := map[string]interface{}{
				"id":            "1",
				"name":          req["name"],
				"metadata":      req["metadata"],
				"retentionDays": int(req["retention"].(float64)),
			}
			json.NewEncoder(w).Encode(resp)
		case http.MethodGet:
			resp := map[string]interface{}{
				"id":            "1",
				"name":          "proj",
				"metadata":      map[string]interface{}{"a": "b"},
				"retentionDays": 5,
			}
			json.NewEncoder(w).Encode(resp)
		case http.MethodPut:
			resp := map[string]interface{}{
				"id":            "1",
				"name":          "proj-upd",
				"metadata":      map[string]interface{}{"a": "b"},
				"retentionDays": 5,
			}
			json.NewEncoder(w).Encode(resp)
		case http.MethodDelete:
			w.WriteHeader(204)
		}
	}))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProject(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("langfuse_project.test", "name", "proj"),
					resource.TestCheckResourceAttr("langfuse_project.test", "retention", "5"),
				),
			},
		},
	})
}

func testAccProject(url string) string {
	return fmt.Sprintf(
		`provider "langfuse" {
  host = "%s"
}

resource "langfuse_project" "test" {
  name = "proj"
  retention = 5
  metadata = {
    a = "b"
  }
}
`, url)
}
