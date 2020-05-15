package github

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-github/v31/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccGithubProjectCard_basic(t *testing.T) {
	var card github.ProjectCard

	rn := "github_project_card.card"
	note := "## Unaccepted :point_down:"
	updatedNote := "## Accepted :point_up:"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGithubProjectCardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubProjectCardConfig(note),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubProjectCardExists(rn, &card),
					testAccCheckGithubProjectCardAttributes(&card, &testAccGithubProjectCardExpectedAttributes{
						Note: note,
					}),
				),
			},
			{
				Config: testAccGithubProjectCardConfig(updatedNote),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubProjectCardExists(rn, &card),
					testAccCheckGithubProjectCardAttributes(&card, &testAccGithubProjectCardExpectedAttributes{
						Note: updatedNote,
					}),
				),
			},
		},
	})
}

func testAccGithubProjectCardDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Organization).v3client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "github_project_card" {
			continue
		}

		cardID, err := strconv.ParseInt(rs.Primary.Attributes["card_id"], 10, 64)
		if err != nil {
			return err
		}

		card, res, err := conn.Projects.GetProjectCard(context.TODO(), cardID)
		if err == nil {
			if card != nil &&
				card.GetID() == cardID {
				return fmt.Errorf("Project Card still exists")
			}
		}
		if res.StatusCode != 404 {
			return err
		}
	}
	return nil
}

func testAccCheckGithubProjectCardExists(n string, card *github.ProjectCard) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		cardID, err := strconv.ParseInt(rs.Primary.Attributes["card_id"], 10, 64)
		if err != nil {
			return err
		}

		conn := testAccProvider.Meta().(*Organization).v3client
		gotCard, _, err := conn.Projects.GetProjectCard(context.TODO(), cardID)
		if err != nil {
			return err
		}
		*card = *gotCard
		return nil
	}
}

type testAccGithubProjectCardExpectedAttributes struct {
	Note string
}

func testAccCheckGithubProjectCardAttributes(card *github.ProjectCard, want *testAccGithubProjectCardExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if note := card.GetNote(); note != want.Note {
			return fmt.Errorf("got project card %q; want %q", note, want.Note)
		}

		return nil
	}
}

func testAccGithubProjectCardConfig(note string) string {
	return fmt.Sprintf(`
resource "github_organization_project" "project" {
  name = "An Organization Project"
  body = "This is an organization project."
}

resource "github_project_column" "column" {
  project_id = github_organization_project.project.id
  name       = "Backlog"
}

resource "github_project_card" "card" {
  column_id   = github_project_column.column.column_id
  note        = "%s"
}
`, note)
}
