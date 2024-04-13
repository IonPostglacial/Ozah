package taxons

import (
	"net/url"

	"nicolas.galipot.net/hazo/server/common"
)

type Panel uint64

const (
	PropertiesPanel = Panel(1 << iota)
	DescriptorsPanel
	SummaryPanel
)

var panelNames = []string{"Properties", "Descriptors", "Summary"}

type PanelSet struct {
	common.BitSet
}

func PanelSetFromString(s string) PanelSet {
	return PanelSet{common.BitSetFromString(s, common.BitSet(PropertiesPanel|DescriptorsPanel|SummaryPanel), common.EmptyBitSet)}
}

func (ps PanelSet) LinkToPanelState(url *url.URL) string {
	query := url.Query()
	query.Del("panels")
	query.Add("panels", ps.String())
	newUrl := *url
	newUrl.RawQuery = query.Encode()
	return newUrl.String()
}
