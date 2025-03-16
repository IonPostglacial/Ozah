package taxons

import (
	"context"
	"fmt"

	"nicolas.galipot.net/hazo/server/action"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/appdb"
)

type menuActions struct {
	cc      *common.Context
	queries *appdb.Queries
}

func (h *menuActions) addMenuLanguage(ctx context.Context, langRef string) error {
	fmt.Printf("add menu language: %s, %s\n", h.cc.User.Login, langRef)
	_, err := h.queries.InsertUserSelectedLanguage(ctx, appdb.InsertUserSelectedLanguageParams{
		UserLogin: h.cc.User.Login,
		LangRef:   langRef,
	})
	return err
}

func (h *menuActions) removeMenuLanguage(ctx context.Context, langRef string) error {
	_, err := h.queries.DeleteUserSelectedLanguage(ctx, appdb.DeleteUserSelectedLanguageParams{
		UserLogin: h.cc.User.Login,
		LangRef:   langRef,
	})
	return err
}

func (h *menuActions) Register(handlers *[]action.Handler) {
	*handlers = append(*handlers, action.NewHandlerWithStringArgument("menu-lang-add", h.addMenuLanguage))
	*handlers = append(*handlers, action.NewHandlerWithStringArgument("menu-lang-remove", h.removeMenuLanguage))
}
