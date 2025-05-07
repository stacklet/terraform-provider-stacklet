package api

type SSOGroupConfig struct {
	Name              string   `json:"name"`
	Roles             []string `json:"roles"`
	AccountGroupUUIDs []string `graphql:"accountGroupUUIDs"`
}
