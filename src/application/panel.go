package application

type Panel uint64

const (
	PropertiesPanel = Panel(iota)
	DescriptorsPanel
	SummaryPanel
)

var PanelNames = []string{"Properties", "Descriptors", "Summary"}

func (p Panel) String() string {
	return PanelNames[p]
}

type UnselectedPanel struct {
	Value uint64
	Name  string
}
