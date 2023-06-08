package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GraphQLQuery struct {
	Query string `json:"query"`
}

func sendGraphQLQuery(query string) {
	url := "https://api.github.com/graphql"
	authToken := "YOUR_GITHUB_AUTH_TOKEN"

	// Create a context with an HTTP client and set headers
	ctx := context.Background()
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(query)))
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// Process the response
	fmt.Println(result)
}

func listRepos(repoName string) (string, error) {
	query := `
    {
        repository(owner: "OWNER_NAME", name: "REPO_NAME") {
            issues(last: 5) {
                nodes {
                    title
                    body
                    url
                }
            }
        }
    }`
	sendGraphQLQuery(query)
}
