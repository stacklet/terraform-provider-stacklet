package provider

import (
	"encoding/base64"
	"encoding/json"

	"github.com/hasura/go-graphql-client"
)

func wrapNodeID(parts []string) graphql.ID {
	jsonBytes, err := json.Marshal(parts)
	if err != nil {
		// This should never happen with a simple string array
		return graphql.ID("")
	}
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return graphql.ID(encoded)
}
