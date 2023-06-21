package main

import (
	"context"

	"github.com/machinebox/graphql"
)

type RepositoryURL struct {
	Repository struct {
		URL string `json:"url"`
	} `json:"repository"`
}

func repoURL(owner string, repoName string) string {
	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	graphqlRequest := graphql.NewRequest(
		`
			query($owner: String!, $repoName: String!) {
				repository(owner: $owner, name: $repoName) {
					url
				}
			}
		`,
	)
	graphqlRequest.Var("owner", owner)
	graphqlRequest.Var("repoName", repoName)

	graphqlRequest.Header.Add("Authorization", "bearer "+getENVVar("GITHUB_TOKEN"))
	var graphqlResponse RepositoryURL

	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}
	return graphqlResponse.Repository.URL
}
