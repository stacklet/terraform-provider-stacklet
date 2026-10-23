// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MaskSecrets returns a context whose logs show the given secrets only as a
// mask. Request bodies carry them JSON-encoded, so that form is masked too.
func MaskSecrets(ctx context.Context, secrets ...*string) context.Context {
	var masked []string
	for _, s := range secrets {
		if s == nil || *s == "" {
			continue
		}
		masked = append(masked, *s)
		if encoded, err := json.Marshal(*s); err == nil {
			if e := string(encoded[1 : len(encoded)-1]); e != *s {
				masked = append(masked, e)
			}
		}
	}
	if len(masked) == 0 {
		return ctx
	}
	return tflog.MaskLogStrings(ctx, masked...)
}
