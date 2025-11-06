// Copyright (c) 2025 - Stacklet, Inc.

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

// RepositoryView is the data returned by reading repository view data.
type RepositoryView struct {
	Namespace         string
	BranchName        string
	PolicyDirectories []string
	PolicyFileSuffix  []string
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

func (a repositoryAPI) Read(ctx context.Context, uuid string) (*Repository, error) {
	var q struct {
		Payload struct {
			RepositoryConfig Repository
			Problems         []Problem
		} `graphql:"repositoryConfig(uuid: $uuid)"`
	}
	if err := a.c.Query(ctx, &q, map[string]any{"uuid": uuid}); err != nil {
		return nil, NewAPIError(err)
	}
	if err := fromProblems(ctx, q.Payload.Problems); err != nil {
		return nil, err
	}
	return &q.Payload.RepositoryConfig, nil
}

func (a repositoryAPI) FindByURL(ctx context.Context, url string) (string, error) {
	cursor := ""
	var q struct {
		Conn struct {
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
	for {
		variables := map[string]any{"cursor": graphql.String(cursor)}
		if err := a.c.Query(ctx, &q, variables); err != nil {
			return "", NewAPIError(err)
		}
		if err := fromProblems(ctx, q.Conn.Problems); err != nil {
			return "", err
		}
		for _, edge := range q.Conn.Edges {
			if edge.Node.URL == url {
				return edge.Node.UUID, nil
			}
		}
		if !q.Conn.PageInfo.HasNextPage {
			return "", NotFound{"Repository not found"}
		}
		cursor = q.Conn.PageInfo.EndCursor
	}
}

func (a repositoryAPI) Create(ctx context.Context, i RepositoryCreateInput) (*Repository, error) {
	var m struct {
		Payload struct {
			RepositoryConfig Repository
			Problems         []Problem
		} `graphql:"addRepositoryConfig(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &m, map[string]any{"input": i}); err != nil {
		return nil, NewAPIError(err)
	}
	if err := fromProblems(ctx, m.Payload.Problems); err != nil {
		return nil, err
	}
	return &m.Payload.RepositoryConfig, nil
}

func (a repositoryAPI) Update(ctx context.Context, i RepositoryUpdateInput) (*Repository, error) {
	var m struct {
		Payload struct {
			RepositoryConfig Repository
			Problems         []Problem
		} `graphql:"updateRepositoryConfig(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &m, map[string]any{"input": i}); err != nil {
		return nil, NewAPIError(err)
	}
	if err := fromProblems(ctx, m.Payload.Problems); err != nil {
		return nil, err
	}
	return &m.Payload.RepositoryConfig, nil
}

func (a repositoryAPI) Delete(ctx context.Context, i RepositoryDeleteInput) error {
	var m struct {
		Payload struct {
			Problems []Problem
		} `graphql:"removeRepositoryConfig(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &m, map[string]any{"input": i}); err != nil {
		return NewAPIError(err)
	}
	if err := fromProblems(ctx, m.Payload.Problems); err != nil {
		return err
	}
	return nil
}
