package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketProject_basic(t *testing.T) {
	var project Project

	testTeam := os.Getenv("BITBUCKET_TEAM")

	testAccBitbucketProjectConfig := fmt.Sprintf(`
		resource "bitbucket_project" "test_project" {
			owner = "%s"
			name = "test-project-for-terraform-test"
			key = "TFTEST"
			is_private = true
			description = "A test project for terraform"
		}
	`, testTeam)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBitbucketProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketProjectExists("bitbucket_project.test_project", &project),
				),
			},
		},
	})
}

func testAccCheckBitbucketProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketClient)
	rs, ok := s.RootModule().Resources["bitbucket_project.test_project"]
	if !ok {
		return fmt.Errorf("Not found %s", "bitbucket_project.test_project")
	}

	response, _ := client.Get(fmt.Sprintf("2.0/teams/%s/projects/%s", rs.Primary.Attributes["owner"], rs.Primary.Attributes["key"]))

	if response.StatusCode != 404 {
		return fmt.Errorf("Project still exists")
	}

	return nil
}

func testAccCheckBitbucketProjectExists(n string, project *Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No project ID is set")
		}
		return nil
	}
}
