package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/MasoudHeydari/eps-api/config"
	"github.com/MasoudHeydari/eps-api/ent"
	_ "github.com/lib/pq"
)

func NewDB(app config.App) (*ent.Client, error) {
	ctx := context.Background()
	client, err := ent.Open(app.DB.Schema, getConnectionURL(app))
	if err != nil {
		return nil, err
	}
	err = client.Schema.Create(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getConnectionURL(app config.App) string {
	u := url.URL{
		Scheme:   app.DB.Schema,
		User:     url.UserPassword(app.DB.User, app.DB.Password),
		Host:     fmt.Sprintf("%s:%d", app.DB.Host, app.DB.Port),
		Path:     app.DB.Name,
		RawQuery: app.DB.RawQuery,
	}
	fmt.Println(u.String())
	return u.String()
}
