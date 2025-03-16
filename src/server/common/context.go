package common

import (
	"database/sql"
	"fmt"
	"html/template"

	"nicolas.galipot.net/hazo/storage"
	"nicolas.galipot.net/hazo/storage/appdb"
	"nicolas.galipot.net/hazo/user"
)

type Context struct {
	User       *user.T
	Template   *template.Template
	Config     *ServerConfig
	appDb      *sql.DB
	appQueries *appdb.Queries
}

func (cc *Context) RegisterUser(login string) error {
	u, err := user.Register(login)
	if err != nil {
		return err
	}
	cc.User = u
	return nil
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
