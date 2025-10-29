// Copyright (c) 2025 - Stacklet, Inc.

package schemavalidate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// UUID is a validator that checks that the value is a valid UUID.
func UUID() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"),
		"must be a valid UUID",
	)
}
