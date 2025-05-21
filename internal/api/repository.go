package api

import (
	"context"
	"encoding/json"

	"github.com/hasura/go-graphql-client"
)

type Repository struct {
	ID          string
	UUID        string
	URL         string
	Name        string
	Description *string
	WebhookURL  string `graphql:"webhookURL"`
	System      bool

	Auth struct {
		AuthUser         *string
		HasAuthToken     bool
		SSHPublicKey     *string
		HasSshPrivateKey bool
		HasSshPassphrase bool
	}
}

// RepositoryConfigAuthInput exists to allow us to set only the fields we want to
// change in this type, which matches the expectations of the API and is much more
// clear to casual inspection than exacting use of `omitempty` struct tags.
type RepositoryConfigAuthInput struct {
	m map[string]any
}

func NewRepositoryConfigAuthInput() *RepositoryConfigAuthInput {
	return &RepositoryConfigAuthInput{m: map[string]any{}}
}

func (i *RepositoryConfigAuthInput) SetAuthUser(v *string) {
	i.m["authUser"] = v
}

func (i *RepositoryConfigAuthInput) SetAuthToken(v *string) {
	i.m["authToken"] = v
}

func (i *RepositoryConfigAuthInput) SetSSHPrivateKey(v *string) {
	i.m["sshPrivateKey"] = v
}

func (i *RepositoryConfigAuthInput) SetSSHPassphrase(v *string) {
	i.m["sshPassphrase"] = v
}

func (i *RepositoryConfigAuthInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.m)
}

func (i *RepositoryConfigAuthInput) GetGraphQLType() string {
	return "RepositoryConfigAuthInput"
}

type RepositoryCreateInput struct {
	URL         string                     `json:"url"`
	Name        string                     `json:"name"`
	Description *string                    `json:"description,omitempty"`
	Auth        *RepositoryConfigAuthInput `json:"auth,omitempty"`
}

func (i RepositoryCreateInput) GetGraphQLType() string {
	return "AddRepositoryConfigInput"
}

type RepositoryUpdateInput struct {
	UUID        string                     `json:"uuid"`
	Name        string                     `json:"name"`
	Description *string                    `json:"description"`
	Auth        *RepositoryConfigAuthInput `json:"auth,omitempty"`
}

func (i RepositoryUpdateInput) GetGraphQLType() string {
	return "UpdateRepositoryConfigInput"
}

type RepositoryDeleteInput struct {
	UUID    string `json:"uuid"`
	Cascade bool   `json:"cascade"`
}

func (i RepositoryDeleteInput) GetGraphQLType() string {
	return "RemoveRepositoryConfigInput"
}

type repositoryAPI struct {
	c *graphql.Client
}

func (a repositoryAPI) Read(ctx context.Context, uuid string) (Repository, error) {
	var query struct {
		RepositoryConfig struct {
			RepositoryConfig Repository
			Problems         []Problem
		} `graphql:"repositoryConfig(uuid: $uuid)"`
	}
	variables := map[string]any{
		"uuid": uuid,
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return Repository{}, APIError{"Client error", err.Error()}
	}
	if err := FromProblems(ctx, query.RepositoryConfig.Problems); err != nil {
		return Repository{}, err
	}
	return query.RepositoryConfig.RepositoryConfig, nil
}

func (a repositoryAPI) FindByURL(ctx context.Context, url string) (string, error) {
	cursor := ""
	for {
		var search struct {
			RepositoryConfigs struct {
				Edges []struct {
					Node struct {
						URL  string
						UUID string
					}
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
				Problems []Problem
			} `graphql:"repositoryConfigs(first: 100, after: $cursor)"`
		}
		if err := a.c.Query(ctx, &search, map[string]any{"cursor": cursor}); err != nil {
			return "", APIError{"Client error", err.Error()}
		}
		connection := search.RepositoryConfigs
		if err := FromProblems(ctx, connection.Problems); err != nil {
			return "", err
		}
		for _, edge := range connection.Edges {
			if edge.Node.URL == url {
				return edge.Node.UUID, nil
			}
		}
		if !connection.PageInfo.HasNextPage {
			return "", NotFound{"Repository with given URL not found"}
		}
		cursor = connection.PageInfo.EndCursor
	}
}

func (a repositoryAPI) Create(ctx context.Context, i RepositoryCreateInput) (Repository, error) {
	var m struct {
		AddRepositoryConfig struct {
			RepositoryConfig Repository
			Problems         []Problem
		} `graphql:"addRepositoryConfig(input: $input)"`
	}
	variables := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &m, variables); err != nil {
		return Repository{}, APIError{"Client error", err.Error()}
	}
	payload := m.AddRepositoryConfig
	if err := FromProblems(ctx, payload.Problems); err != nil {
		return Repository{}, err
	}
	return payload.RepositoryConfig, nil
}

func (a repositoryAPI) Update(ctx context.Context, i RepositoryUpdateInput) (Repository, error) {
	var m struct {
		UpdateRepositoryConfig struct {
			RepositoryConfig Repository
			Problems         []Problem
		} `graphql:"updateRepositoryConfig(input: $input)"`
	}
	variables := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &m, variables); err != nil {
		return Repository{}, APIError{"Client error", err.Error()}
	}
	payload := m.UpdateRepositoryConfig
	if err := FromProblems(ctx, payload.Problems); err != nil {
		return Repository{}, err
	}
	return payload.RepositoryConfig, nil
}

func (a repositoryAPI) Delete(ctx context.Context, i RepositoryDeleteInput) error {
	input := map[string]any{"input": i}
	var m struct {
		RemoveRepositoryConfig struct {
			Problems []Problem
		} `graphql:"removeRepositoryConfig(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &m, input); err != nil {
		return APIError{"Client error", err.Error()}
	}
	if err := FromProblems(ctx, m.RemoveRepositoryConfig.Problems); err != nil {
		return err
	}
	return nil
}
