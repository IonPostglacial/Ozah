package authentication

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/common"
)

const SessionCookieName = "session_token"

//go:embed login.html
var loginTemplate string

type Model struct {
	ErrorMessage string
}

func HandlerWrapper(handler common.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		cdb, queries, err := db.OpenCommon()
		if err != nil {
			log.Fatal(err)
		}
		token, err := r.Cookie(SessionCookieName)
		if err == nil {
			username, err := loginFromSessionToken(ctx, queries, token.Value)
			if err != nil {
				fmt.Printf("error logging: %s\n", err)
				// TODO: handle case with error not login error
			} else {
				handler(w, r, &common.Context{
					User: common.User{
						Login: username,
					},
				})
				return
			}
		} else {
			fmt.Printf("error reading cookie: %s\n", err)
		}
		loginFound := false
		err = r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		username := r.Form.Get("login")
		password := r.Form.Get("password")
		cred, err := queries.GetCredentials(ctx, username)
		if err != nil {
			fmt.Printf("login not found: %s\n", err)
			// TODO: login not found
		} else {
			loginFound = true
		}
		authorized := false
		if loginFound {
			switch cred.Encryption {
			case "bcrypt":
				err := bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(password))
				authorized = (err == nil)
			}
		}
		if authorized {
			session, err := startSession(ctx, cdb, queries, username)
			if err != nil {
				log.Fatal(err)
			}
			http.SetCookie(w, &http.Cookie{
				Name:     SessionCookieName,
				Value:    session.Token,
				Expires:  session.Expires,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
			})
			handler(w, r, &common.Context{
				User: common.User{
					Login: username,
				},
			})
		} else {
			tmpl := template.New("login")
			template.Must(tmpl.Parse(loginTemplate))
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			if loginFound {
				err := tmpl.Execute(w, &Model{
					ErrorMessage: fmt.Sprintf("Sorry '%s': wrong password", username),
				})
				if err != nil {
					log.Fatal(err)
				}
			} else {
				msg := ""
				if username != "" {
					msg = fmt.Sprintf("Login '%s': not found", username)
				}
				err := tmpl.Execute(w, &Model{
					ErrorMessage: msg,
				})
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
