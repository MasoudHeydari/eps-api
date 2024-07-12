package db

import (
	"context"

	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/MasoudHeydari/eps-api/ent/serp"
	"github.com/MasoudHeydari/eps-api/model"
)

func InsertSearchResult(ctx context.Context, db *ent.Client, result model.SearchResult, sqID int) error {
	err := db.SERP.Create().
		SetTitle(result.Title).
		SetDescription(result.Description).
		SetURL(result.URL).
		SetKeyWords(nil2Zero(result.KeyWords)).
		SetEmails(nil2Zero(result.Emails)).
		SetPhones(nil2Zero(result.Phones)).
		SetSqID(sqID).
		OnConflictColumns(
			serp.FieldSqID,
			serp.FieldURL,
			serp.FieldEmails,
			serp.FieldPhones,
		).
		UpdateNewValues().
		Exec(ctx)
	return err
}

func nil2Zero(s []string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	return s
}
