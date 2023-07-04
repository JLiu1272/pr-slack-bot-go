package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func sendResponseToChannel(w http.ResponseWriter, r *http.Request) {

	// Unmarshal the request body
	var reqBody slack.InteractionCallback
	json.Unmarshal([]byte(r.FormValue("payload")), &reqBody)

	// print the unmarsheled request body
	// print the unmarsheled request body with formatting. Where the key is indented and the value is on the same line
	fmt.Printf("Request Body: %+v\n\n", reqBody)

	w.Write([]byte(helpMessege()))
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	command := s.Command
	params := &slack.Msg{Text: s.Text}

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
			response := interface{}(nil)

			if len(actions) == 3 {
				username := actions[2]
				response = listAction(repoName, username)
			} else {
				response = listAction(repoName, s.UserName)
			}

			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				// Handle error
				http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
				return
			}
			return
		}

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func main() {
	http.HandleFunc("/receive", slashCommandHandler)
	http.HandleFunc("/send", sendResponseToChannel)

	// log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func main() {
// 	block := build_block("JLiu1272", "github-webhook-server")

// 	jsonBytes, err := json.MarshalIndent(block, "", " ")
// 	if err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return
// 	}

// 	jsonString := string(jsonBytes)
// 	fmt.Println(jsonString)
// }
