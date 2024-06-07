package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"Todo.com/m/pkg/models"
)

var Data []*models.Todo

// tasks is a handler
// Allows: GET, POST
// if GET - Show the list of tasks
// if POST - adds user task
func (app *application) tasks(w http.ResponseWriter, r *http.Request) {
	//parsing all the file
	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")

	//if POST POST appedend the data
	if r.Method == http.MethodPost {
		details := r.FormValue("task")
		tags := r.FormValue("tags")
		if len(strings.TrimSpace(details)) != 0 {
			_, err := app.todo.Insert(details, tags)
			if err != nil {
				app.errorLog.Fatal(err.Error())
			}
			app.session.Put(r, "flash", "Task successfully created!")
		} else {
			app.session.Put(r, "flash", "Task Creation Failed")
		}
	}

	s, err := app.todo.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

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
		http.Error(w, "Internal Server Error", http.StatusNotImplemented)
	}

}

// Func Delete is a handler
// recieves ID from the client Request
// Delete the data based on the id
// Redirect to the main page
func (app *application) deleteTask(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	errDel, isSpecial := app.todo.Delete(title)
	fmt.Println(isSpecial)
	if isSpecial {
		fmt.Println(r.FormValue("title"))
		app.special.Delete(r.FormValue("title"))
	}

	if errDel != nil {
		app.errorLog.Println(errDel)
		app.session.Put(r, "flash", "Task cannot be Deleted!")
	} else {
		app.session.Put(r, "flash", "Task successfully Deleted!")
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	idToUpdate, _ := strconv.Atoi(r.FormValue("id"))
	details := r.FormValue("update")
	if len(strings.TrimSpace(details)) != 0 {
		errUpd := app.todo.Update(idToUpdate, r.FormValue("update"))
		if errUpd != nil {
			app.errorLog.Println(errUpd)
		} else {
			app.session.Put(r, "flash", "Task successfully Updated!")
		}
	} else {
		app.session.Put(r, "flash", "Task cannot be Updated!")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/login.page.tmpl")
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 501)
	} else {
		err = ts.Execute(w, app.session.PopString(r, "flash"))
	}
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username, password)

	//check if user exist:
	//if yes, then return true else false
	isUser, err := app.user.IsThereAnUser(username, password)

	if err != nil {
		app.errorLog.Println(err)
	}
	if isUser {
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "flash", "Login Successful")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		app.session.Put(r, "Authenticated", false)
		app.session.Put(r, "flash", "Wrong Username or Password")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

}

func (app *application) signUpForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/signup.page.tmpl")
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 501)
	} else {
		err = ts.Execute(w, app.session.PopString(r, "flash"))
	}
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username, password)

	//check if user exist:
	//if yes, then return true else false
	err := app.user.SignUp(username, password)

	if err != nil {
		app.errorLog.Println(err)
		app.session.Put(r, "flash", "User Already Exist")
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
	} else {
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "flash", "SignUp Successful")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.session.PopBool(r, "Authenticated")
	app.session.Put(r, "flash", "User Logged Out")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) tags(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/tags.page.tmpl")
	tag := r.FormValue("t")
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusNotImplemented)
	}
	data, todoerr := app.todo.GetAll()
	if todoerr != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusNotImplemented)
	}

	Tags := []models.Tag{}

	for _, v := range data {
		for _, i := range v.Tags {
			if i == tag {
				t := models.Tag{
					Name: tag,
					Task: *v,
				}
				Tags = append(Tags, t)
			}
		}
	}

	err = ts.Execute(w, Tags) //send tag data
	if err != nil {
		app.errorLog.Println(err)
	}

}
