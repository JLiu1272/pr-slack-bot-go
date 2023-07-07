package main

type Block struct {
	Type     string     `json:"type"`
	Text     *TextBlock `json:"text,omitempty"`
	Elements []Element  `json:"elements,omitempty"`
	Style    string     `json:"style,omitempty"`
	Indent   int        `json:"indent,omitempty"`
}

type TextBlock struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

type Element struct {
	Type        string     `json:"type"`
	URL         string     `json:"url,omitempty"`
	Text        string     `json:"text,omitempty"`
	Elements    []Element  `json:"elements,omitempty"`
	Style       string     `json:"style,omitempty"`
	Indent      int        `json:"indent,omitempty"`
	Options     []Option   `json:"options,omitempty"`
	Placeholder *TextBlock `json:"placeholder,omitempty"`
	ActionID    string     `json:"action_id,omitempty"`
}

type Option struct {
	Text  *TextBlock `json:"text,omitempty"`
	Value string     `json:"value,omitempty"`
}

type Blocks struct {
	Blocks []Block `json:"blocks"`
}

func richTextSection(repoInfo Repository, repoName string) []Element {
	rich_text_section := []Element{}
	for _, pr := range repoInfo.Repository.PullRequests.Nodes {
		rich_text_section = append(rich_text_section, Element{
			Type: "rich_text_section",
			Elements: []Element{
				{
					Type: "link",
					URL:  pr.URL,
					Text: pr.Title,
				},
			},
		})
	}

	return rich_text_section
}

func selectOptions(repoInfo Repository, repoName string) []Option {
	select_options := []Option{}

	for _, pr := range repoInfo.Repository.PullRequests.Nodes {
		select_options = append(select_options, Option{
			Text: &TextBlock{
				Type:  "plain_text",
				Text:  pr.Title,
				Emoji: true,
			},
			Value: pr.URL,
		})
	}

	return select_options
}

func listPRsMessage(username string, repoName string) (prsMessage interface{}) {

	prs, err := listPRs(username, repoName)

	if err != nil {
		panic(err)
	}

	rich_text_section := richTextSection(prs, repoName)
	select_options := selectOptions(prs, repoName)

	prsMessage = Blocks{
		Blocks: []Block{
			{
				Type: "section",
				Text: &TextBlock{
					Type: "mrkdwn",
					Text: "Here are the top 5 PRs:",
				},
			},
			{
				Type: "rich_text",
				Elements: []Element{
					{
						Type:     "rich_text_list",
						Elements: rich_text_section,
						Style:    "bullet",
						Indent:   0,
					},
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "actions",
				Elements: []Element{
					{
						Type: "static_select",
						Placeholder: &TextBlock{
							Type:  "plain_text",
							Text:  "Select a PR to Send",
							Emoji: true,
						},
						Options:  select_options,
						ActionID: "select-reminder-time",
					},
				},
			},
		},
	}

	return prsMessage
}
