package common

import "html/template"

type User struct {
	Login string
}

type Context struct {
	User     *User
	Template *template.Template
}
