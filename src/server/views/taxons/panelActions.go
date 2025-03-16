package taxons

import (
	"context"

	"nicolas.galipot.net/hazo/application"
	"nicolas.galipot.net/hazo/server/action"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/appdb"
)

type panelActions struct {
	cc      *common.Context
	queries *appdb.Queries
}

func (h *panelActions) addPanel(ctx context.Context, panelId int) error {
	_, err := h.queries.DeleteUserHiddenPanels(ctx, appdb.DeleteUserHiddenPanelsParams{
		UserLogin: h.cc.User.Login,
		PanelID:   int64(panelId),
	})
	return err
}

func (h *panelActions) removePanel(ctx context.Context, panelId int) error {
	_, err := h.queries.InsertUserHiddenPanels(ctx, appdb.InsertUserHiddenPanelsParams{
		UserLogin: h.cc.User.Login,
		PanelID:   int64(panelId),
	})
	return err
}

func (h *panelActions) zoomPanel(ctx context.Context, panelId int) error {
	for id := range application.PanelNames {
		_, err := h.queries.InsertUserHiddenPanels(ctx, appdb.InsertUserHiddenPanelsParams{
			UserLogin: h.cc.User.Login,
			PanelID:   int64(id),
		})
		if err != nil {
			return err
		}
	}
	_, err := h.queries.DeleteUserHiddenPanels(ctx, appdb.DeleteUserHiddenPanelsParams{
		UserLogin: h.cc.User.Login,
		PanelID:   int64(panelId),
	})
	return err
}

func (h *panelActions) Register(handlers *[]action.Handler) {
	*handlers = append(*handlers, action.NewHandlerWithIntArgument("panel-add", h.addPanel))
	*handlers = append(*handlers, action.NewHandlerWithIntArgument("panel-remove", h.removePanel))
	*handlers = append(*handlers, action.NewHandlerWithIntArgument("panel-zoom", h.zoomPanel))
}
