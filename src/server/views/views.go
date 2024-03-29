package views

import (
	"fmt"

	"nicolas.galipot.net/hazo/server/components/popover"
)

func NewMenuState(label, dsName string) *popover.State {
	return &popover.State{
		Label: label,
		Items: []popover.Item{
			{Url: fmt.Sprintf("/ds/%s/taxons", dsName), Label: "Taxons"},
			{Url: fmt.Sprintf("/ds/%s/characters", dsName), Label: "Characters"},
		},
	}
}
