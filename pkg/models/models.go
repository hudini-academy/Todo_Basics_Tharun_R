package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Todo struct {
	Id      int
	Title   string
	Created time.Time
	Expires time.Time
	Tags    []string
}

type User struct {
	Id       int
	Username string
	Password string
}

type Tag struct {
	Name string
	Task Todo
}
