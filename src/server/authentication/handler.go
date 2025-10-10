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
		token, err := r.Cookie(SessionCookieName)
		if err == nil {
			username, err := loginFromSessionToken(ctx, cc.AppQueries(), token.Value)
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
		cred, err := cc.AppQueries().GetCredentials(ctx, username)
		if err == nil {
			loginFound = true
		}
		authorized := false
		isMSAccount := false
		if loginFound {
			switch cred.Encryption {
			case "bcrypt":
				err := bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(password))
				authorized = (err == nil)
			case "ms-oauth":
				isMSAccount = true
			}
		}
		if authorized {
			session, err := startSession(ctx, cc, username)
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
			if isMSAccount {
				err := tmpl.Execute(w, &Model{
					ErrorMessage: fmt.Sprintf("Account '%s' uses Microsoft authentication. Please use the 'Sign in with Microsoft' button.", username),
				})
				if err != nil {
					return err
				}
			} else if loginFound {
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

func MSLoginHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	config := LoadMSConfig()
	HandleMSLogin(w, r, config)
	return nil
}

func MSCallbackHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	config := LoadMSConfig()
	HandleMSCallback(w, r, config, cc.AppQueries())
	return nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	ctx := context.Background()

	token, err := r.Cookie(SessionCookieName)
	if err == nil {
		cc.AppQueries().DeleteSession(ctx, token.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
