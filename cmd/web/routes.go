// handling the routes
package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

// Function to handle all the routes for the mux
func (app *application) routes() http.Handler {
	mux := pat.New()
	//signup and logins
	mux.Get("/signup", app.session.Enable(http.HandlerFunc(app.signUpForm)))
	mux.Post("/signup", app.session.Enable(http.HandlerFunc(app.signUp)))
	mux.Post("/login", app.session.Enable(http.HandlerFunc(app.login)))
	mux.Get("/login", app.session.Enable(http.HandlerFunc(app.loginPage)))
	mux.Get("/logout", app.session.Enable(http.HandlerFunc(app.logout)))
	mux.Get("/", app.session.Enable(app.authenticate(http.HandlerFunc(app.tasks))))
	//special route
	mux.Get("/special", app.session.Enable(app.authenticate(http.HandlerFunc(app.special_tasks))))
	mux.Post("/special/delete", app.session.Enable(http.HandlerFunc(app.special_deleteTask)))
	mux.Post("/special/update", app.session.Enable(http.HandlerFunc(app.special_updateTask)))
	//Todo route
	mux.Post("/tasks", app.session.Enable(app.authenticate(http.HandlerFunc(app.tasks))))
	mux.Post("/delete", app.session.Enable(http.HandlerFunc(app.deleteTask)))
	mux.Post("/update", app.session.Enable(http.HandlerFunc(app.updateTask)))
	//Tags
	mux.Get("/tags", app.session.Enable(http.HandlerFunc(app.tags)))
	//Panic
	mux.Get("/panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Intentional panic for testing purposes")
	}))
	//adding the static css file to fileserver
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
