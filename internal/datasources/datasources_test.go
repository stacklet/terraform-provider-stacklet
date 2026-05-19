// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatasourcesList_ReleasedOnly(t *testing.T) {
	d := datasources{
		Released: []func() datasource.DataSource{
			newFactory(&accountDataSource{}),
			newFactory(&policyDataSource{}),
		},
		Unreleased: []func() datasource.DataSource{
			newFactory(&roleDataSource{}),
		},
	}

	result := d.List(false)

	require.Len(t, result, 2)
	assert.IsType(t, &accountDataSource{}, result[0]())
	assert.IsType(t, &policyDataSource{}, result[1]())
}

func TestDatasourcesList_IncludeUnreleased(t *testing.T) {
	d := datasources{
		Released: []func() datasource.DataSource{
			newFactory(&accountDataSource{}),
			newFactory(&policyDataSource{}),
		},
		Unreleased: []func() datasource.DataSource{
			newFactory(&roleDataSource{}),
		},
	}

	result := d.List(true)

	require.Len(t, result, 3)
	assert.IsType(t, &accountDataSource{}, result[0]())
	assert.IsType(t, &policyDataSource{}, result[1]())
	assert.IsType(t, &roleDataSource{}, result[2]())
}
