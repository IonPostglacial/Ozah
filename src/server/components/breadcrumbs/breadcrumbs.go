package breadcrumbs

type BreadCrumb struct {
	Label string
	Url   string
}

type State struct {
	Branch []BreadCrumb
}
