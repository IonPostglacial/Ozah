package authentication

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/storage"
)

const SessionCookieName = "session_token"

//go:embed login.html
var loginTemplate string

type Model struct {
	ErrorMessage string
}

func HandlerWrapper(handler common.Handler) common.Handler {
	return func(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
		ctx := context.Background()
		cdb, queries, err := storage.OpenAppDb()
		if err != nil {
			return err
		}
		token, err := r.Cookie(SessionCookieName)
		if err == nil {
			username, err := loginFromSessionToken(ctx, queries, token.Value)
			if err != nil {
				fmt.Printf("error logging: %s\n", err)
				// TODO: handle case with error not login error
			} else {
				if err := cc.RegisterUser(username); err != nil {
					return err
				}
				return handler(w, r, cc)
			}
		}
		loginFound := false
		err = r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		username := r.Form.Get("login")
		password := r.Form.Get("password")
		cred, err := queries.GetCredentials(ctx, username)
		if err == nil {
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
			if err := cc.RegisterUser(username); err != nil {
				return err
			}
			return handler(w, r, cc)
		} else {
			tmpl := components.NewTemplate()
			template.Must(tmpl.Parse(loginTemplate))
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			if loginFound {
				err := tmpl.Execute(w, &Model{
					ErrorMessage: fmt.Sprintf("Sorry '%s': wrong password", username),
				})
				if err != nil {
					return err
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
					return err
				}
			}
		}
		return nil
	}
}
