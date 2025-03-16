package authentication

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/appdb"
)

type Session struct {
	Token   string
	Expires time.Time
}

var ErrSessionExpired = fmt.Errorf("session expired")

const dateFormat = "2006-01-02T15:04:05"

func loginFromSessionToken(ctx context.Context, queries *appdb.Queries, tok string) (string, error) {
	session, err := queries.GetSession(ctx, tok)
	if err != nil {
		return "", err
	}
	expiryDate, err := time.Parse(dateFormat, session.ExpiryDate)
	if err != nil {
		return "", err
	}
	if expiryDate.Before(time.Now()) {
		return "", ErrSessionExpired
	}
	return session.Login, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var tokenSize = 512

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

func startSession(ctx context.Context, cc *common.Context, login string) (*Session, error) {
	sessionToken, err := generateToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(2 * time.Hour)
	err = cc.AppQueriesTx(func(qtx *appdb.Queries) error {
		_, err = qtx.DeleteUserSessions(ctx, login)
		if err != nil {
			return err
		}
		_, err = qtx.InsertSession(ctx, appdb.InsertSessionParams{
			Token:      sessionToken,
			Login:      login,
			ExpiryDate: expiresAt.Format(dateFormat),
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Session{Token: sessionToken, Expires: expiresAt}, nil
}
