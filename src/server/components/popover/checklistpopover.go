package popover

type CheckListItem struct {
	Checked     bool
	ActionName  string
	ActionValue string
	Label       string
}

type CheckListViewModel struct {
	Label string
	Icon  string
	Items []CheckListItem
}
