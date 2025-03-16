package documents

import (
	"context"
	"fmt"

	"nicolas.galipot.net/hazo/server/action"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/appdb"
)

type MenuActions struct {
	cc *common.Context
}

func NewMenuActions(cc *common.Context) *MenuActions {
	return &MenuActions{cc}
}

func (h *MenuActions) addMenuLanguage(ctx context.Context, langRef string) error {
	fmt.Printf("add menu language: %s, %s\n", h.cc.User.Login, langRef)
	_, err := h.cc.AppQueries().InsertUserSelectedLanguage(ctx, appdb.InsertUserSelectedLanguageParams{
		UserLogin: h.cc.User.Login,
		LangRef:   langRef,
	})
	return err
}

func (h *MenuActions) removeMenuLanguage(ctx context.Context, langRef string) error {
	_, err := h.cc.AppQueries().DeleteUserSelectedLanguage(ctx, appdb.DeleteUserSelectedLanguageParams{
		UserLogin: h.cc.User.Login,
		LangRef:   langRef,
	})
	return err
}

func (h *MenuActions) Register(reg *action.Registry) {
	reg.AppendAction(action.NewActionWithStringArgument("menu-lang-add", h.addMenuLanguage))
	reg.AppendAction(action.NewActionWithStringArgument("menu-lang-remove", h.removeMenuLanguage))
}
