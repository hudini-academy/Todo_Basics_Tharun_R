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
}

type User struct {
	Id       int
	Username string
	Password string
}
