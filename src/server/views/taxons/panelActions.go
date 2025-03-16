package taxons

import (
	"context"

	"nicolas.galipot.net/hazo/application"
	"nicolas.galipot.net/hazo/server/action"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/appdb"
)

type panelActions struct {
	cc *common.Context
}

func NewPanelActions(cc *common.Context) *panelActions {
	return &panelActions{cc}
}

func (h *panelActions) addPanel(ctx context.Context, panelId int) error {
	_, err := h.cc.AppQueries().DeleteUserHiddenPanels(ctx, appdb.DeleteUserHiddenPanelsParams{
		UserLogin: h.cc.User.Login,
		PanelID:   int64(panelId),
	})
	return err
}

func (h *panelActions) removePanel(ctx context.Context, panelId int) error {
	_, err := h.cc.AppQueries().InsertUserHiddenPanels(ctx, appdb.InsertUserHiddenPanelsParams{
		UserLogin: h.cc.User.Login,
		PanelID:   int64(panelId),
	})
	return err
}

func (h *panelActions) zoomPanel(ctx context.Context, panelId int) error {
	for id := range application.PanelNames {
		_, err := h.cc.AppQueries().InsertUserHiddenPanels(ctx, appdb.InsertUserHiddenPanelsParams{
			UserLogin: h.cc.User.Login,
			PanelID:   int64(id),
		})
		if err != nil {
			return err
		}
	}
	_, err := h.cc.AppQueries().DeleteUserHiddenPanels(ctx, appdb.DeleteUserHiddenPanelsParams{
		UserLogin: h.cc.User.Login,
		PanelID:   int64(panelId),
	})
	return err
}

func (h *panelActions) Register(reg *action.Registry) {
	reg.AppendAction(action.NewActionWithIntArgument("panel-add", h.addPanel))
	reg.AppendAction(action.NewActionWithIntArgument("panel-remove", h.removePanel))
	reg.AppendAction(action.NewActionWithIntArgument("panel-zoom", h.zoomPanel))
}
