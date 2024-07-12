package model

type APIResponse struct {
	Version       string  `json:"version"`
	StatusCode    int     `json:"status_code"`
	StatusMessage string  `json:"status_message"`
	Time          string  `json:"time"`
	Cost          float64 `json:"cost"`
	TasksCount    int     `json:"tasks_count"`
	TasksError    int     `json:"tasks_error"`
	Tasks         []Task  `json:"tasks"`
}

type Task struct {
	ID            string   `json:"id"`
	StatusCode    int      `json:"status_code"`
	StatusMessage string   `json:"status_message"`
	Time          string   `json:"time"`
	Cost          float64  `json:"cost"`
	ResultCount   int      `json:"result_count"`
	Path          []string `json:"path"`
	Data          TaskData `json:"data"`
	Result        []Result `json:"result"`
}

type TaskData struct {
	API          string `json:"api"`
	Function     string `json:"function"`
	SE           string `json:"se"`
	SEType       string `json:"se_type"`
	Depth        int    `json:"depth"`
	Keyword      string `json:"keyword"`
	LanguageCode string `json:"language_code"`
	LocationCode int    `json:"location_code"`
	Device       string `json:"device"`
	OS           string `json:"os"`
}

type Result struct {
	Keyword      string `json:"keyword"`
	Type         string `json:"type"`
	SEDomain     string `json:"se_domain"`
	LocationCode int    `json:"location_code"`
	LanguageCode string `json:"language_code"`
	CheckURL     string `json:"check_url"`
	Datetime     string `json:"datetime"`
	// Spell          string   `json:"spell"`
	ItemTypes      []string `json:"item_types"`
	SEResultsCount int      `json:"se_results_count"`
	ItemsCount     int      `json:"items_count"`
	Items          []Item   `json:"items"`
}

type Item struct {
	Type         string `json:"type"`
	RankGroup    int    `json:"rank_group"`
	RankAbsolute int    `json:"rank_absolute"`
	Domain       string `json:"domain"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	URL          string `json:"url"`
	Breadcrumb   string `json:"breadcrumb"`
}
