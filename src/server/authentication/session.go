package authentication

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"time"

	"nicolas.galipot.net/hazo/db/commonstorage"
)

type Session struct {
	Token   string
	Expires time.Time
}

var ErrSessionExpired = fmt.Errorf("session expired")

const dateFormat = "2006-01-02T15:04:05"

func loginFromSessionToken(ctx context.Context, queries *commonstorage.Queries, tok string) (string, error) {
	session, err := queries.GetSession(ctx, tok)
	if err != nil {
		return "", err
	}
	expiryDate, err := time.Parse(dateFormat, session.ExpiryDate)
	if err != nil {
		return "", err
	}
	fmt.Printf("expiry date: %s\nnow: %s\n", expiryDate, time.Now())
	if expiryDate.Before(time.Now()) {
		return "", ErrSessionExpired
	}
	return session.Login, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var tokenSize = 512

// n is the length of random string we want to generate
func generateToken() (string, error) {
	var buf strings.Builder

	for range tokenSize {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters)-1)))
		if err != nil {
			return "", err
		}
		buf.WriteRune(letters[n.Int64()])
	}
	return buf.String(), nil
}

func startSession(ctx context.Context, db *sql.DB, queries *commonstorage.Queries, login string) (*Session, error) {
	sessionToken, err := generateToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(2 * time.Hour)
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	_, err = qtx.DeleteUserSessions(ctx, login)
	if err != nil {
		return nil, err
	}
	_, err = qtx.InsertSession(ctx, commonstorage.InsertSessionParams{
		Token:      sessionToken,
		Login:      login,
		ExpiryDate: expiresAt.Format(dateFormat),
	})
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return &Session{Token: sessionToken, Expires: expiresAt}, nil
}
