package langfuse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPromptDataSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		resp := map[string]interface{}{
			"version":       2,
			"type":          "prompt",
			"prompt":        map[string]interface{}{"role": "system"},
			"config":        map[string]interface{}{"a": "b"},
			"labels":        []string{"l1"},
			"tags":          []string{"t1"},
			"commitMessage": "msg",
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
				Config: testAccDataSourcePrompt(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.langfuse_prompt.test", "type", "prompt"),
					resource.TestCheckResourceAttr("data.langfuse_prompt.test", "version_out", "2"),
					resource.TestCheckResourceAttr("data.langfuse_prompt.test", "labels.0", "l1"),
					resource.TestCheckResourceAttr("data.langfuse_prompt.test", "tags.0", "t1"),
					resource.TestCheckResourceAttr("data.langfuse_prompt.test", "commit_message", "msg"),
				),
			},
		},
	})
}

func testAccDataSourcePrompt(url string) string {
	return fmt.Sprintf(`
provider "langfuse" {
  host = "%s"
}

data "langfuse_prompt" "test" {
  name = "prompt"
}
`, url)
}
