package common

import (
	"html/template"

	"nicolas.galipot.net/hazo/user"
)

type Context struct {
	User     *user.T
	Template *template.Template
	Config   *ServerConfig
}

func (cc *Context) RegisterUser(login string) error {
	u, err := user.Register(login)
	if err != nil {
		return err
	}
	cc.User = u
	return nil
}
