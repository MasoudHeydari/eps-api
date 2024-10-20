package config

type App struct {
	Version        string   `json:"version"`
	Http           HTTP     `json:"http"`
	Limiter        Limiter  `json:"limiter"`
	DB             Database `json:"database"`
	QueryDepth     int      `json:"query_depth"`
	FileNameMaxLen int      `json:"file_name_max_len"`
	APIKey         string   `json:"api_key"`
}

type HTTP struct {
	Listen string `json:"listen"`
	Port   string `json:"port"`
}

type Database struct {
	Schema   string `json:"schema"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	RawQuery string `json:"raw_query"`
}

type Limiter struct {
	Hour   uint `json:"hour"`
	Minute uint `json:"minute"`
	Burst  int  `json:"burst"`
}

func New() App {
	return App{}
}
