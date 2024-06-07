package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"Todo.com/m/pkg/models"
)

var Data1 []*models.Todo

// tasks is a handler
// Allows: GET, POST
// if GET - Show the list of tasks
// if POST - adds user task
func (app *application) special_tasks(w http.ResponseWriter, r *http.Request) {
	//parsing all the file
	ts, err := template.ParseFiles("./ui/html/special.page.tmpl")

	//if POST POST appedend the data
	s, err := app.special.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	Data = s
	err = ts.Execute(w, struct {
		Tasks []*models.Todo
		Flash string
	}{
		Tasks: s,
		Flash: app.session.PopString(r, "flash"),
	})
	//send output to the client

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 501)
	}

}

// Func Delete is a handler
// recieves ID from the client Request
// Delete the data based on the id
// Redirect to the main page
func (app *application) special_deleteTask(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	errDel := app.special.Delete(title)
	IstodoTitle, _ := app.todo.IsThereATodo(title)
	if IstodoTitle {
		app.todo.Delete(title)
	}
	if errDel != nil {
		app.errorLog.Println(errDel)
		app.session.Put(r, "flash", "Task cannot be Deleted!")
	} else {
		app.session.Put(r, "flash", "Task successfully Deleted!")
	}

	// Redirect to home page
	http.Redirect(w, r, "/special", http.StatusSeeOther)
}

func (app *application) special_updateTask(w http.ResponseWriter, r *http.Request) {
	idToUpdate, _ := strconv.Atoi(r.FormValue("id"))
	details := r.FormValue("update")
	if len(strings.TrimSpace(details)) != 0 {
		errUpd := app.special.Update(idToUpdate, r.FormValue("update"))
		if errUpd != nil {
			app.errorLog.Println(errUpd)
		} else {
			app.session.Put(r, "flash", "Task successfully Updated!")
		}
	} else {
		app.session.Put(r, "flash", "Task cannot be Updated!")
	}

	http.Redirect(w, r, "/special", http.StatusSeeOther)
}
