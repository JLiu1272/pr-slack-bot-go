package main

import (
	"context"

	"github.com/machinebox/graphql"
)

type Repository struct {
	Repository Issues `json:"repository"`
}

type Issues struct {
	Issues Nodes `json:"issues"`
}

type Nodes struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

func listRepos(owner string, repoName string) Repository {
	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	graphqlRequest := graphql.NewRequest(
		`
			query($owner: String!, $repoName: String!) {
				repository(owner: $owner, name: $repoName) {
					issues(last: 5) {
						nodes {
							title
							body
							url
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
		panic(err)
	}
	return graphqlResponse
}
