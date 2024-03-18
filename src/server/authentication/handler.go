package authentication

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	w.Header().Add("WWW-Authenticate", "Basic")
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf("<!DOCTYPE html><html>%d %s<br>%s<br>%s<br>%v</html>",
		http.StatusUnauthorized,
		http.StatusText(http.StatusUnauthorized),
		username, password, ok)))
}
