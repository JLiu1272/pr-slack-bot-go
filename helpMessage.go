package main

func helpMessege() string {
	response := "Usage: /pr <command> [arguments]\n"
	response += "Commands:\n"
	response += "  list <repo-name> - List the 5 most recent pull requests for the specified repository\n"
	response += "  list <repo-name> <github-username> - List the top 5 most recent PRs in that repo that was opened by the specified user\n"
	response += "  help - Display this help message\n"
	return response
}
