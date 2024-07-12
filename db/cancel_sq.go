package db

import (
	"context"
	"github.com/MasoudHeydari/eps-api/ent"
)

func CancelSQ(ctx context.Context, db *ent.Client, sqID int) error {
	_, err := db.SearchQuery.UpdateOneID(sqID).
		SetIsFinished(true).
		Save(ctx)
	return err
}
