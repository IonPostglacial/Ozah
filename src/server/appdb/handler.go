package appdb

import (
	"net/http"

	"nicolas.galipot.net/hazo/server/common"
)

func Handler(handler common.Handler) common.Handler {
	return func(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
		err := cc.ConnectAppDb()
		if err != nil {
			return err
		}
		return handler(w, r, cc)
	}
}
