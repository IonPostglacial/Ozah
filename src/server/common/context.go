package common

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"nicolas.galipot.net/hazo/server/action"
	"nicolas.galipot.net/hazo/storage"
	"nicolas.galipot.net/hazo/storage/appdb"
	"nicolas.galipot.net/hazo/user"
)

type Context struct {
	User           *user.T
	Template       *template.Template
	Config         *ServerConfig
	appDb          *sql.DB
	appQueries     *appdb.Queries
	actionRegistry *action.Registry
}

func NewContext(config *ServerConfig) *Context {
	return &Context{Config: config, actionRegistry: action.NewRegistry()}
}

func (cc *Context) RegisterUser(login string) error {
	u, err := user.Register(login)
	if err != nil {
		return err
	}
	cc.User = u
	return nil
}

func (cc *Context) RegisterActions(reg action.Registrable) {
	cc.actionRegistry.Register(reg)
}

func (cc *Context) ExecuteActions(ctx context.Context, r *http.Request) error {
	return cc.actionRegistry.ExecuteActions(ctx, r)
}

func (cc *Context) ConnectAppDb() error {
	db, appQueries, err := storage.OpenAppDb()
	if err != nil {
		return fmt.Errorf("couldn't open global database: %w", err)
	}
	cc.appDb = db
	cc.appQueries = appQueries
	return err
}

func (cc *Context) AppQueries() *appdb.Queries {
	return cc.appQueries
}

func (cc *Context) AppQueriesTx(cb func(*appdb.Queries) error) error {
	tx, err := cc.appDb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := cc.appQueries.WithTx(tx)
	cb(qtx)
	return tx.Commit()
}
