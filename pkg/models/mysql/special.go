package mysql

import (
	"database/sql"

	"Todo.com/m/pkg/models"
)

// Define a TodoModel type which wraps a sql.DB connection pool.
type SpTodoModel struct {
	DB *sql.DB
}

func (t *SpTodoModel) Insert(title string) (int, error) {
	stmt := `INSERT INTO special_task (title, created, expires)
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

func (t *SpTodoModel) Delete(title string) error {
	stmt := `Delete from special_task where title = ?`

	_, err := t.DB.Exec(stmt, title)
	if err != nil {
		return err
	}
	return nil
}

func (m *SpTodoModel) GetAll() ([]*models.Todo, error) {
	stmt := `SELECT * from special_task`
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

func (m *SpTodoModel) Update(id int, title string) error {
	stmt := `Update special_task set title =? WHERE sid=?`
	_, err := m.DB.Exec(stmt, title, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *SpTodoModel) SearchTile(title string) (bool, error) {
	stmt := `SELECT * from todos where title=?`
	rows, err := m.DB.Query(stmt, title)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
