package popover

type Item struct {
	Url   string
	Label string
}

type State struct {
	Label string
	Items []Item
}
