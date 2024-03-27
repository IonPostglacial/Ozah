package components

import (
	"net/http"

	"nicolas.galipot.net/hazo/server/common"
)

func JavascriptHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	w.Header().Add("Content-Type", "application/javascript")
	w.Write([]byte(JavascriptCode))
	return nil
}
