//go:build generate

package tools

import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)

// Generate documentation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .. // generate-docs

// Validate generated documentation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate --provider-dir .. // validate-docs
