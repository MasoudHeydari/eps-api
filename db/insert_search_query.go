package db

import (
	"context"

	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/MasoudHeydari/eps-api/ent/searchquery"
	"github.com/google/uuid"
)

func InsertNewSearchQuery(ctx context.Context, db *ent.Client, loc int, lang, searchQ string) (int, error) {
	sqEnt, err := db.SearchQuery.Create().
		SetLocCode(loc).
		SetLanguage(lang).
		SetQuery(searchQ).
		Save(ctx)
	if err != nil {
		switch {
		case ent.IsConstraintError(err):
			err = db.SearchQuery.Update().
				Where(searchquery.IsFinished(true)).
				SetIsFinished(false).
				Exec(ctx)
			if err != nil {
				return -1, err
			}
		default:
			return -1, err
		}
	}
	return sqEnt.ID, nil
}

func InsertJobID(ctx context.Context, db *ent.Client, sqID int, jobID uuid.UUID) error {
	_, err := db.SearchQuery.Update().
		Where(searchquery.ID(sqID)).
		SetJobID(jobID).
		Save(ctx)
	return err
}

func GetAJobID(ctx context.Context, db *ent.Client) (uuid.UUID, int, error) {
	i, err := db.SearchQuery.Query().
		Where(searchquery.IsFinished(false)).
		Order(searchquery.ByID()).
		First(ctx)
	if err != nil {
		return uuid.Nil, -1, err
	}
	return i.JobID, i.ID, nil
}
