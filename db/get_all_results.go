package db

import (
	"context"
	"fmt"
	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/MasoudHeydari/eps-api/ent/serp"
	"github.com/MasoudHeydari/eps-api/model"
)

func GetAllResult(ctx context.Context, db *ent.Client, sqID, page int) ([]model.SERP, error) {
	offSet := page * 5
	entSERPs, err := db.SERP.Query().
		Where(
			serp.SqID(sqID),
			serp.IsRead(false),
		).
		Offset(offSet).
		Limit(5).
		Order(ent.Desc(serp.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]model.SERP, 0, len(entSERPs))
	for _, entSERP := range entSERPs {
		fmt.Println("kw: ", entSERP.KeyWords)
		results = append(results,
			model.SERP{
				URL:         entSERP.URL,
				Title:       entSERP.Title,
				Description: entSERP.Description,
				Phones:      entSERP.Phones,
				Emails:      entSERP.Emails,
				Keywords:    entSERP.KeyWords,
				IsRead:      entSERP.IsRead,
				CreatedAt:   entSERP.CreatedAt,
			},
		)
	}
	return results, nil
}
