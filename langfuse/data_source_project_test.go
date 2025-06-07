package langfuse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		resp := map[string]interface{}{
			"id":            "1",
			"name":          "proj",
			"metadata":      map[string]interface{}{"a": "b"},
			"retentionDays": 5,
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
				Config: testAccDataSourceProject(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.langfuse_project.test", "name", "proj"),
					resource.TestCheckResourceAttr("data.langfuse_project.test", "retention", "5"),
					resource.TestCheckResourceAttr("data.langfuse_project.test", "metadata.a", "b"),
				),
			},
		},
	})
}

func testAccDataSourceProject(url string) string {
	return fmt.Sprintf(`
provider "langfuse" {
  host = "%s"
}

data "langfuse_project" "test" {
  id = "1"
}
`, url)
}
