package common

import "net/http"

type Handler func(http.ResponseWriter, *http.Request, *Context)
