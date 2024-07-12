package delivery

type searchQ struct {
	Language string `json:"lang"`
	LocCode  int    `json:"loc"`
	Query    string `json:"q"`
}

type CancelSQ struct {
	SQID int `json:"sq_id"` // SQID is CreateJob Query ID
}

type GetAllSearchResults struct {
	SQID int `param:"sq_id"` // SQID is CreateJob Query ID
	Page int `query:"page"`
}

type ExportCSV struct {
	SQID int `param:"sq_id"` // SQID is CreateJob Query ID
}
