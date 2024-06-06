package mysql

import (
	"database/sql"

	"Todo.com/m/pkg/models"
)

// Define a TodoModel type which wraps a sql.DB connection pool.
type TodoModel struct {
	DB *sql.DB
}

func (t *TodoModel) Insert(title string) (int, error) {
	stmt := `INSERT INTO todos (title, created, expires)
    VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := t.DB.Exec(stmt, title, 365)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (t *TodoModel) Delete(id int) error {
	stmt := `Delete from todos where id = ?`
	_, err := t.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
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
		err = rows.Scan(&s.Id, &s.Title, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of todo.
		todo = append(todo, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the todo slice.
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
