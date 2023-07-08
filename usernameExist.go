package main

import (
	"context"

	"github.com/machinebox/graphql"
)

type Username struct {
	User struct {
		Login string `json:"login"`
	} `json:"user"`
}

func usernameExist(username string) bool {
	client := graphql.NewClient("https://api.github.com/graphql")
	request := graphql.NewRequest(
		`
			query($username: String!) {
				user(login: $username) {
					login
				}	
			}
		`,
	)
	request.Var("username", username)

	request.Header.Add("Authorization", "bearer "+getENVVar("GITHUB_TOKEN"))

	resp := &Username{}
	if err := client.Run(context.Background(), request, &resp); err != nil {
		return false
	}

	return true
}
