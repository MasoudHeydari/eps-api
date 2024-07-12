package model

type SearchResult struct {
	Rank        int      `json:"rank"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Phones      []string `json:"phones"`
	Emails      []string `json:"emails"`
	KeyWords    []string `json:"key_words"`
	Description string   `json:"description"`
}
