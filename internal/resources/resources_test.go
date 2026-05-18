// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourcesList_ReleasedOnly(t *testing.T) {
	r := resources{
		Released: []func() resource.Resource{
			newFactory(&accountResource{}),
			newFactory(&policyCollectionResource{}),
		},
		Unreleased: []func() resource.Resource{
			newFactory(&roleAssignmentResource{}),
		},
	}

	result := r.List(false)

	require.Len(t, result, 2)
	assert.IsType(t, &accountResource{}, result[0]())
	assert.IsType(t, &policyCollectionResource{}, result[1]())
}

func TestResourcesList_IncludeUnreleased(t *testing.T) {
	r := resources{
		Released: []func() resource.Resource{
			newFactory(&accountResource{}),
			newFactory(&policyCollectionResource{}),
		},
		Unreleased: []func() resource.Resource{
			newFactory(&roleAssignmentResource{}),
		},
	}

	result := r.List(true)

	require.Len(t, result, 3)
	assert.IsType(t, &accountResource{}, result[0]())
	assert.IsType(t, &policyCollectionResource{}, result[1]())
	assert.IsType(t, &roleAssignmentResource{}, result[2]())
}
