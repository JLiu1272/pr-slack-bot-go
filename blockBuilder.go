package main

type Block struct {
	Blocks []BlockElement `json:"blocks"`
}

type BlockElement struct {
	Type     string    `json:"type"`
	Text     *Text     `json:"text,omitempty"`
	Elements []Element `json:"elements,omitempty"`
}

type Text struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

type Element struct {
	Type     string `json:"type"`
	Text     Text   `json:"text"`
	Value    string `json:"value"`
	ActionID string `json:"action_id"`
}

func build_block() (data interface{}) {

	data_typed := Block{
		Blocks: []BlockElement{
			{
				Type: "section",
				Text: &Text{
					Type: "mrkdwn",
					Text: "New Paid Time Off request from <example.com|Fred Enriquez>\n\n<https://example.com|View request>",
				},
			},
			{
				Type: "actions",
				Elements: []Element{
					{
						Type: "button",
						Text: Text{
							Type:  "plain_text",
							Text:  "Send PR to Chat",
							Emoji: true,
						},
						Value:    "click_me_123",
						ActionID: "actionId-0",
					},
				},
			},
		},
	}

	var data_interface interface{} = data_typed

	return data_interface
}
