package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

func listRepos() {
	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	graphqlRequest := graphql.NewRequest(
		`
			{
				repository(owner: "JLiu1272", name: "github-webhook-server") {
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
	graphqlRequest.Header.Add("Authorization", "bearer "+getENVVar("GITHUB_TOKEN"))
	var graphqlResponse interface{}
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}
	fmt.Println(graphqlResponse)
}
