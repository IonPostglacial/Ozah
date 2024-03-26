package common

import (
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request, *Context) error

func UnwrapHandler(handler Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r, &Context{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
