// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/stretchr/testify/assert"
)

func TestMaskSecrets(t *testing.T) {
	plain := "plain-secret"
	escaped := "escaped\"secret"
	empty := ""
	encoded, _ := json.Marshal(map[string]string{"key": escaped})

	var output bytes.Buffer
	ctx := tflogtest.RootLogger(context.Background(), &output)
	ctx = MaskSecrets(ctx, &plain, &escaped, &empty, nil)
	tflog.Debug(ctx, "message", map[string]any{"plain": plain, "body": string(encoded)})

	entries, err := tflogtest.MultilineJSONDecode(&output)
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "***", entries[0]["plain"])
	assert.Equal(t, `{"key":"***"}`, entries[0]["body"])
	assert.NotContains(t, fmt.Sprint(entries[0]), "secret")
}
