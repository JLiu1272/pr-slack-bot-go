package main

import "fmt"

func listAction(repoName string, userName string) interface{} {
	if !usernameExist(userName) {
		return fmt.Sprintf("Username: %v does not exist. Please provide the username using `/pr list [repo-name] [github-username]` command", userName)
	}

	return listPRsMessage(userName, repoName)
}
