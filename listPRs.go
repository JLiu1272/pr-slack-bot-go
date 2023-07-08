package main

import (
	"context"

	"github.com/machinebox/graphql"
)

type Repository struct {
	Repository PullRequests `json:"repository"`
}

type PullRequests struct {
	PullRequests Nodes `json:"pullRequests"`
}

type Nodes struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	CreatedAt string `json:"createdAt"`
}

func listPRs(owner string, repoName string) (Repository, error) {
	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	graphqlRequest := graphql.NewRequest(
		`
			query($owner: String!, $repoName: String!) {
				repository(owner: $owner, name: $repoName) {
					pullRequests(last: 5, orderBy: {field: CREATED_AT, direction: DESC}) {
						nodes {
							title
							url
							author {
								login
							}
							createdAt
						}
					}
				}
			}
		`,
	)
	graphqlRequest.Var("owner", owner)
	graphqlRequest.Var("repoName", repoName)

	graphqlRequest.Header.Add("Authorization", "bearer "+getENVVar("GITHUB_TOKEN"))
	var graphqlResponse Repository

	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return graphqlResponse, err
	}
	return graphqlResponse, nil
}
