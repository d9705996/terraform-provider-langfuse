package langfuse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccProjectResource(t *testing.T) {
	var deleted bool
	var updated bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			resp := map[string]interface{}{
				"id":            "1",
				"name":          req["name"],
				"metadata":      req["metadata"],
				"retentionDays": int(req["retention"].(float64)),
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("encode response: %v", err)
			}
		case http.MethodGet:
			name := "proj"
			if updated {
				name = "proj-upd"
			}
			resp := map[string]interface{}{
				"id":            "1",
				"name":          name,
				"metadata":      map[string]interface{}{"a": "b"},
				"retentionDays": 5,
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("encode response: %v", err)
			}
		case http.MethodPut:
			resp := map[string]interface{}{
				"id":            "1",
				"name":          "proj-upd",
				"metadata":      map[string]interface{}{"a": "b"},
				"retentionDays": 5,
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("encode response: %v", err)
			}
			updated = true
		case http.MethodDelete:
			deleted = true
			w.WriteHeader(204)
		}
	}))
	defer server.Close()

	name1 := "proj"
	name2 := "proj-upd"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(s *terraform.State) error {
			if !deleted {
				return fmt.Errorf("project not deleted")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccProject(server.URL, name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("langfuse_project.test", "name", name1),
					resource.TestCheckResourceAttr("langfuse_project.test", "retention", "5"),
				),
			},
			{
				Config: testAccProjectUpdated(server.URL, name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("langfuse_project.test", "name", name2),
				),
			},
		},
	})
}

func testAccProject(url, name string) string {
	return fmt.Sprintf(`
provider "langfuse" {
  host = "%s"
}

resource "langfuse_project" "test" {
  name = "%s"
  retention = 5
  metadata = {
    a = "b"
  }
}
`, url, name)
}

func testAccProjectUpdated(url, name string) string {
	return fmt.Sprintf(`
provider "langfuse" {
  host = "%s"
}

resource "langfuse_project" "test" {
  name = "%s"
  retention = 5
  metadata = {
    a = "b"
  }
}
`, url, name)
}
