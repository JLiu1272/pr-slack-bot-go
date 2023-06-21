package main

import "fmt"

func listAction(repoName string, userName string) string {
	if !usernameExist(userName) {
		return fmt.Sprintf("Username: %v does not exist. Please provide the username using `/pr list <repo-name> <github-username>` command", userName)
	}

	prs, err := listPRs(userName, repoName)

	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	response := formatListPRsResponse(prs, repoName)
	return response

}
