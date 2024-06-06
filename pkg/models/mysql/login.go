package mysql

import (
	"database/sql"
	"fmt"
)

// Define a TodoModel type which wraps a sql.DB connection pool.
type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) IsThereAnUser(username, password string) (bool, error) {
	stmt := "select id from users where username = ? and password = ?"
	rows, err := u.DB.Query(stmt, username, password)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (u *UserModel) SignUp(username, password string) error {
	stmt := "insert into users (username, password) values( ?, ?)"
	result, err := u.DB.Exec(stmt, username, password)
	fmt.Println("mdoel", result)
	if err != nil {
		return err
	}
	return nil
}
