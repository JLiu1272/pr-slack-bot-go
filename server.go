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

func hash_str(message string, secretKey string) string {
	// Convert the secret key to bytes
	key := []byte(secretKey)

	// Create an HMAC-SHA256 hasher
	hasher := hmac.New(sha256.New, key)

	// Write the message to the hasher
	hasher.Write([]byte(message))

	// Get the final hash
	hash := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash)

	fmt.Println("HMAC SHA256 hash:", hashString)
	return hashString
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
	generated_hash := hash_str(baseStr, getENVVar("OLD_TOKEN"))

	fmt.Printf("X-Slack-Signature: %v\n\n", slack_signature)
	fmt.Printf("baseStr: %v\n\n", baseStr)
	fmt.Printf("Generated Hash: %v\n\n", generated_hash)

	// if generated_hash != slack_signature {
	// 	fmt.Println("Invalid token")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	switch command {
	case "/pr":
		response := fmt.Sprintf("You asked for the weather for %v", params.Text)
		w.Write([]byte(response))

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
