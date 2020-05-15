---
layout: "github"
page_title: "GitHub: github_project_card"
description: |-
  Creates and manages project cards for GitHub projects
---

# github_project_card

This resource allows you to create and manage cards for GitHub projects.

## Example Usage

```hcl
resource "github_organization_project" "project" {
  name = "An Organization Project"
  body = "This is an organization project."
}

resource "github_project_column" "column" {
  project_id = github_organization_project.project.id
  name       = "Backlog"
}

resource "github_project_card" "card" {
  project_id  = github_organization_project.project.id
  column_id = github_organization_project.column.column_id
  note        = "## Unaccepted 👇"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of an existing project that the card will be created in.

* `column_id` - (Required) The ID of the card.

* `note` - (Required) The note contents of the card. Markdown supported.
