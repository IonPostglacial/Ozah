package common

type User struct {
	Login string
}

type Context struct {
	User *User
}
