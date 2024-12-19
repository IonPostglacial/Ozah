package components

import (
	"net/http"

	"nicolas.galipot.net/hazo/server/common"
)

func CssHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	w.Header().Add("Content-Type", "text/css")
	w.Write([]byte(CssCode))
	return nil
}
