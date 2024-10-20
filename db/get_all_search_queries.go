package db

import (
	"context"

	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/MasoudHeydari/eps-api/ent/searchquery"
	"github.com/MasoudHeydari/eps-api/model"
)

func GetAllSearchQueries(ctx context.Context, db *ent.Client, page int) ([]model.SearchQuery, error) {
	offSet := page * 5
	entSearchQueries, err := db.SearchQuery.Query().
		Offset(offSet).
		Limit(5).
		Order(ent.Desc(searchquery.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	searchQueries := make([]model.SearchQuery, 0, len(entSearchQueries))
	for _, entSearchQuery := range entSearchQueries {
		searchQueries = append(searchQueries, model.SearchQuery{
			Id:         entSearchQuery.ID,
			Query:      entSearchQuery.Query,
			Language:   entSearchQuery.Language,
			Location:   entSearchQuery.LocCode,
			IsCanceled: entSearchQuery.IsFinished,
			CreatedAt:  entSearchQuery.CreatedAt,
		})
	}
	return searchQueries, nil
}
