package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"Todo.com/m/pkg/models"
)

// Define a TodoModel type which wraps a sql.DB connection pool.
type TodoModel struct {
	DB *sql.DB
}

func (t *TodoModel) checkIfSpecial(title string) bool {
	if len(title) < 8 {
		return false
	}
	subset := strings.ToLower(title[:8])
	return subset == "special:"

}

func (t *TodoModel) Insert(title, tags string) (int, error) {

	stmt := `INSERT INTO todos (title, created, expires, tags)
    VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`
	result, err := t.DB.Exec(stmt, title, 365, tags)

	if t.checkIfSpecial(title) {
		sp_stmt := `INSERT INTO special_task (title, created, expires)
		VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
		_, err := t.DB.Exec(sp_stmt, title, 365)
		if err != nil {
			return 0, nil
		}
		_, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	}

	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Println("New Record added")
	return int(id), nil
}

func (t *TodoModel) Delete(title string) (error, bool) {
	stmt := `Delete from todos where title = ?`
	//titlesp, _ := t.SearchTile(title)
	_, err := t.DB.Exec(stmt, title)
	if err != nil {
		return err, false
	}
	return nil, t.checkIfSpecial(title)
}

func (m *TodoModel) GetAll() ([]*models.Todo, error) {
	stmt := `SELECT * from todos`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the models.ToDo objects.
	todo := []*models.Todo{}
	for rows.Next() {
		// Create a pointer to a new zeroed ToDo struct.
		s := &models.Todo{}
		tmp := ""
		err = rows.Scan(&s.Id, &s.Title, &s.Created, &s.Expires, &tmp)
		s.Tags = strings.Split(tmp, ",")
		if err != nil {
			fmt.Println(err, 80)
			return nil, err
		}
		// Append it to the slice of todo.
		todo = append(todo, s)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err, 87)
		return nil, err
	}
	// If everything went OK then return the todo slice.
	fmt.Println("Data Fetch Successful")
	return todo, nil
}

func (m *TodoModel) Update(id int, title string) error {
	stmt := `Update todos set title =? WHERE id=?`
	_, err := m.DB.Exec(stmt, title, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *TodoModel) IsThereATodo(title string) (bool, error) {
	stmt := `SELECT * from todos where title=?`
	rows, err := m.DB.Query(stmt, title)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (m *TodoModel) GetTags(tag string) (bool, error) {
	stmt := `select * from todos where tags like "%?%"`
	rows, _ := m.DB.Query(stmt, tag)
	defer rows.Close()
	return rows.Next(), nil
}
