package common

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request, *Context) error
type HandlerWrapper func(Handler) Handler

type Model struct {
	ErrorMessage string
}

//go:embed error.html
var errorPage string

func (handler Handler) Unwrap(config *ServerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.New("error")
		tmpl = template.Must(tmpl.Parse(errorPage))
		err := handler(w, r, NewContext(config))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tmpl.Execute(w, &Model{ErrorMessage: err.Error()})
		}
	}
}

func (h Handler) Wrap(wrapper HandlerWrapper) Handler {
	return wrapper(h)
}
