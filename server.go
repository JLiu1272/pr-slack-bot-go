package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

type BlockActionsEvent struct {
	Actions []ActionBlock `json:"actions"`
}

type ActionBlock struct {
	Type           string         `json:"type"`
	ActionID       string         `json:"action_id"`
	BlockID        string         `json:"block_id"`
	SelectedOption SelectedOption `json:"selected_option"`
	Placeholder    TextBlock      `json:"placeholder"`
	ActionTS       string         `json:"action_ts"`
}

type SelectedOption struct {
	Text  TextBlock `json:"text"`
	Value string    `json:"value"`
}

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

func sendMessageToChannel(message string) (channelID string, timestamp string, err error) {
	api := slack.New(getENVVar("SLACK_BOT_TOKEN"))

	channelID, timestamp, err = api.PostMessage(
		getENVVar("CHANNEL_ID"),
		slack.MsgOptionText(message, false),
	)

	if err != nil {
		return channelID, timestamp, err
	}
	return channelID, timestamp, nil
}

func receiveSelectedOption(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message BlockActionsEvent

	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prTitle := message.Actions[0].SelectedOption.Text.Text
	prURL := message.Actions[0].SelectedOption.Value

	requestCodeReviewMsg := fmt.Sprintf("Can I get a code review request for PR: <%v|%v>", prTitle, prURL)

	w.Write([]byte("ACK Message Sent\n"))

	if channelID, timestamp, err := sendMessageToChannel(requestCodeReviewMsg); err != nil {
		log.Printf("[ERROR] Failed to send message to channel: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.Write([]byte(fmt.Sprintf("Message successfully sent to channel %s at %s", channelID, timestamp)))
	}
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
	http.HandleFunc("/send", receiveSelectedOption)

	// log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
