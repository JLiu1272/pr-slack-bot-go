package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

func getENVVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

func formatSlashCommand(command slack.SlashCommand) string {
	var keys []string
	values := make(url.Values)

	// Manually set the order of keys
	keys = append(keys, "token", "team_id", "team_domain", "channel_id", "channel_name", "user_id", "user_name", "command", "text", "response_url", "trigger_id")

	// Set the values for each key
	values.Set("token", command.Token)
	values.Set("team_id", command.TeamID)
	values.Set("team_domain", command.TeamDomain)
	values.Set("channel_id", command.ChannelID)
	values.Set("channel_name", command.ChannelName)
	values.Set("user_id", command.UserID)
	values.Set("user_name", command.UserName)
	values.Set("command", command.Command)
	values.Set("text", command.Text)
	values.Set("response_url", command.ResponseURL)
	values.Set("trigger_id", command.TriggerID)

	// Encode the values in the desired order
	var buf bytes.Buffer
	for _, key := range keys {
		if vs := values[key]; len(vs) > 0 {
			for _, v := range vs {
				buf.WriteString(key)
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(v))
				if key != "trigger_id" {
					buf.WriteByte('&')
				}
			}
		}
	}
	return buf.String()

}

func VerifySlackRequest(r *http.Request, baseStr string, signingSecret string) (bool, error) {
	signature := r.Header.Get("X-Slack-Signature")

	// Calculate HMAC-SHA256 hash
	sig := hmac.New(sha256.New, []byte(signingSecret))
	sig.Write([]byte(baseStr))
	expectedSignature := hex.EncodeToString(sig.Sum(nil))

	fmt.Printf("Expected Signature: %v\n\n", expectedSignature)
	fmt.Printf("Slack Signature: %v\n\n", signature)

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature)), nil
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	command := s.Command
	params := &slack.Msg{Text: s.Text}

	slack_signature := r.Header.Get("X-Slack-Signature")
	baseStr := fmt.Sprintf("v0:%v:%v",
		r.Header.Get("X-Slack-Request-Timestamp"),
		formatSlashCommand(s),
	)

	fmt.Printf("X-Slack-Signature: %v\n\n", slack_signature)
	fmt.Printf("baseStr: %v\n\n", baseStr)

	isVerified, err := VerifySlackRequest(r, baseStr, getENVVar("SLACK_TOKEN"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("It can here")
		return
	}

	fmt.Printf("slack request is valid:%v", isVerified)

	switch command {
	case "/pr":
		actions := strings.Fields(params.Text)

		if len(actions) < 1 {
			w.Write([]byte(helpMessege()))
			return
		}

		if actions[0] == "help" {
			w.Write([]byte(helpMessege()))
			return
		}

		if len(actions) < 2 {
			w.Write([]byte(helpMessege()))
			return
		}

		if actions[0] == "list" {
			repoName := actions[1]
			response := ""

			if len(actions) == 3 {
				username := actions[2]
				response = listAction(repoName, username)
			} else {
				response = listAction(repoName, s.UserName)
			}

			w.Write([]byte(response))
			return
		}

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func main() {
	http.HandleFunc("/receive", slashCommandHandler)

	// log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func main() {
// 	fmt.Printf("Username: %v\n", usernameExist("saurabh0719"))
// 	fmt.Printf("List Actions: %v\n", listAction("github-webhook-server", "saurabh0719"))
// }
