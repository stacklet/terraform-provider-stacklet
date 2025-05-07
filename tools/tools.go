//go:build tools

// Run "go generate" to format example terraform files and generate the docs for the registry/website
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

package tools

import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
