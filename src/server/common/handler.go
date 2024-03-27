package common

import (
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request, *Context) error
type HandlerWrapper func(Handler) Handler

func (handler Handler) Unwrap() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r, &Context{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h Handler) Wrap(wrapper HandlerWrapper) Handler {
	return wrapper(h)
}
