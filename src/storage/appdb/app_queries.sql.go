// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: app_queries.sql

package appdb

import (
	"context"
	"database/sql"
)

const deleteUserHiddenPanels = `-- name: DeleteUserHiddenPanels :execresult
delete from User_Hidden_Panel
where
    User_Login = ? and Panel_Id = ?
`

type DeleteUserHiddenPanelsParams struct {
	UserLogin string
	PanelID   int64
}

func (q *Queries) DeleteUserHiddenPanels(ctx context.Context, arg DeleteUserHiddenPanelsParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteUserHiddenPanels, arg.UserLogin, arg.PanelID)
}

const deleteUserSelectedLanguage = `-- name: DeleteUserSelectedLanguage :execresult
delete from User_Selected_Lang
where
    User_Login = ?
    and Lang_Ref = ?
`

type DeleteUserSelectedLanguageParams struct {
	UserLogin string
	LangRef   string
}

func (q *Queries) DeleteUserSelectedLanguage(ctx context.Context, arg DeleteUserSelectedLanguageParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteUserSelectedLanguage, arg.UserLogin, arg.LangRef)
}

const deleteUserSessions = `-- name: DeleteUserSessions :execresult
delete from Session
where
    Login = ?
`

func (q *Queries) DeleteUserSessions(ctx context.Context, login string) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteUserSessions, login)
}

const getCredentials = `-- name: GetCredentials :one
select
    Encryption,
    Password,
    Created_On,
    Last_Modified
from
    Credentials
where
    Login = ?
`

type GetCredentialsRow struct {
	Encryption   string
	Password     string
	CreatedOn    sql.NullString
	LastModified sql.NullString
}

func (q *Queries) GetCredentials(ctx context.Context, login string) (GetCredentialsRow, error) {
	row := q.db.QueryRowContext(ctx, getCredentials, login)
	var i GetCredentialsRow
	err := row.Scan(
		&i.Encryption,
		&i.Password,
		&i.CreatedOn,
		&i.LastModified,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
select
    Login,
    Expiry_Date
from
    Session
where
    Token = ?
`

type GetSessionRow struct {
	Login      string
	ExpiryDate string
}

func (q *Queries) GetSession(ctx context.Context, token string) (GetSessionRow, error) {
	row := q.db.QueryRowContext(ctx, getSession, token)
	var i GetSessionRow
	err := row.Scan(&i.Login, &i.ExpiryDate)
	return i, err
}

const getUserConfiguration = `-- name: GetUserConfiguration :one
select
    login, private_directory
from
    User_Configuration
where
    Login = ?
`

func (q *Queries) GetUserConfiguration(ctx context.Context, login string) (UserConfiguration, error) {
	row := q.db.QueryRowContext(ctx, getUserConfiguration, login)
	var i UserConfiguration
	err := row.Scan(&i.Login, &i.PrivateDirectory)
	return i, err
}

const getUserHiddenPanels = `-- name: GetUserHiddenPanels :many
select
    Panel_Id
from
    User_Hidden_Panel
where
    User_Login = ?
`

func (q *Queries) GetUserHiddenPanels(ctx context.Context, userLogin string) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getUserHiddenPanels, userLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var panel_id int64
		if err := rows.Scan(&panel_id); err != nil {
			return nil, err
		}
		items = append(items, panel_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserSelectedLanguages = `-- name: GetUserSelectedLanguages :many
select
    Lang_Ref
from
    User_Selected_Lang
where
    User_Login = ?
`

func (q *Queries) GetUserSelectedLanguages(ctx context.Context, userLogin string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getUserSelectedLanguages, userLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var lang_ref string
		if err := rows.Scan(&lang_ref); err != nil {
			return nil, err
		}
		items = append(items, lang_ref)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertCredentials = `-- name: InsertCredentials :execresult
insert into
    Credentials (Login, Encryption, Password)
values
    (?, ?, ?)
`

type InsertCredentialsParams struct {
	Login      string
	Encryption string
	Password   string
}

func (q *Queries) InsertCredentials(ctx context.Context, arg InsertCredentialsParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertCredentials, arg.Login, arg.Encryption, arg.Password)
}

const insertLang = `-- name: InsertLang :execresult
insert into
    Lang (Ref, Name)
values
    (?, ?)
`

type InsertLangParams struct {
	Ref  string
	Name string
}

func (q *Queries) InsertLang(ctx context.Context, arg InsertLangParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertLang, arg.Ref, arg.Name)
}

const insertSession = `-- name: InsertSession :execresult
insert into
    Session (Token, Login, Expiry_Date)
values
    (?, ?, ?)
`

type InsertSessionParams struct {
	Token      string
	Login      string
	ExpiryDate string
}

func (q *Queries) InsertSession(ctx context.Context, arg InsertSessionParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertSession, arg.Token, arg.Login, arg.ExpiryDate)
}

const insertUserConfiguration = `-- name: InsertUserConfiguration :execresult
insert into
    User_Configuration (Login, Private_Directory)
values
    (?, ?)
`

type InsertUserConfigurationParams struct {
	Login            string
	PrivateDirectory string
}

func (q *Queries) InsertUserConfiguration(ctx context.Context, arg InsertUserConfigurationParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertUserConfiguration, arg.Login, arg.PrivateDirectory)
}

const insertUserHiddenPanels = `-- name: InsertUserHiddenPanels :execresult
insert into
    User_Hidden_Panel (User_Login, Panel_Id)
values
    (?, ?)
`

type InsertUserHiddenPanelsParams struct {
	UserLogin string
	PanelID   int64
}

func (q *Queries) InsertUserHiddenPanels(ctx context.Context, arg InsertUserHiddenPanelsParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertUserHiddenPanels, arg.UserLogin, arg.PanelID)
}

const insertUserPanel = `-- name: InsertUserPanel :execresult
insert into
    Panel (Id, Name)
values
    (?, ?)
`

type InsertUserPanelParams struct {
	ID   int64
	Name string
}

func (q *Queries) InsertUserPanel(ctx context.Context, arg InsertUserPanelParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertUserPanel, arg.ID, arg.Name)
}

const insertUserSelectedLanguage = `-- name: InsertUserSelectedLanguage :execresult
insert into
    User_Selected_Lang (User_Login, Lang_Ref)
values
    (?, ?)
`

type InsertUserSelectedLanguageParams struct {
	UserLogin string
	LangRef   string
}

func (q *Queries) InsertUserSelectedLanguage(ctx context.Context, arg InsertUserSelectedLanguageParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertUserSelectedLanguage, arg.UserLogin, arg.LangRef)
}
