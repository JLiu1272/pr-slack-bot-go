package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func listRepos() {
	jsonData := map[string]string{
		"query": `
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
			}`,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header.Set("Authorization", "bearer "+getENVVar("GITHUB_TOKEN"))

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))
}
