package breadcrumbs

type BreadCrumb struct {
	Label string
	Url   string
}

type ViewModel struct {
	Branch []BreadCrumb
}
