package popover

type Item struct {
	Url   string
	Label string
}

type ViewModel struct {
	Label string
	Items []Item
}
