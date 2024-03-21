package authentication

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"nicolas.galipot.net/hazo/db"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()

	ctx := context.Background()
	queries, err := db.OpenCommon()
	if err != nil {
		log.Fatal(err)
	}
	loginFound := false
	cred, err := queries.GetCredentials(ctx, username)
	if err != nil {
		fmt.Println("login not found")
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
		http.ServeContent(w, r, "index.html", time.Now(), strings.NewReader("<!DOCTYPE html><html><body>Success</body></html>"))
	} else {
		w.Header().Add("WWW-Authenticate", "Basic")
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnauthorized)
		if loginFound {
			w.Write([]byte(fmt.Sprintf("<!DOCTYPE html><html><body>Sorry '%s': wrong password: '%s'</body></html>",
				username, password)))
		} else {
			w.Write([]byte(fmt.Sprintf("<!DOCTYPE html><html><body>Login '%s': not found</body></html>", username)))
		}
	}
}
