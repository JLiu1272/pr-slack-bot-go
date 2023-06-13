package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

func listRepos(owner string, repoName string) {
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
	var graphqlResponse interface{}
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}
	fmt.Println(graphqlResponse)
}
