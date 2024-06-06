package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"Todo.com/m/pkg/models/mysql"
	"github.com/bmizerany/pat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// Struct for dependency Injection
type application struct {
	todo     *mysql.TodoModel
	user     *mysql.UserModel
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:root@/todo?parseTime=True", "MySql Database String")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key session")
	flag.Parse()

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	//Accesing a file for storing the Information Logs
	infoF, errI := os.OpenFile("./tmp/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errI != nil {
		log.Fatal(errI)
	}
	defer infoF.Close()

	//Accessing a file for storing the Error Logs
	errF, errE := os.OpenFile("./tmp/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errE != nil {
		log.Fatal(errE)
	}
	defer errF.Close()

	//creating a custom logggers
	infoLog := log.New(infoF, "INFO\t", log.LstdFlags)
	errorLog := log.New(errF, "ERROR\t", log.LstdFlags|log.Lshortfile)

	//Connecting to database
	db, sqlerr := openDB(*dsn)
	if sqlerr != nil {
		errorLog.Fatal(sqlerr)
	} else {
		infoLog.Println("connection is established")
	}

	defer db.Close()

	app := &application{
		todo:     &mysql.TodoModel{DB: db},
		user:     &mysql.UserModel{DB: db},
		infoLog:  infoLog,
		errorLog: errorLog,
		session:  session,
	}

	//handling the routes
	mux := pat.New()
	mux.Get("/signup", app.session.Enable(http.HandlerFunc(app.signUpForm)))
	mux.Post("/signup", app.session.Enable(http.HandlerFunc(app.signUp)))
	mux.Post("/login", app.session.Enable(http.HandlerFunc(app.login)))
	mux.Get("/login", app.session.Enable(http.HandlerFunc(app.loginPage)))
	mux.Get("/logout", app.session.Enable(http.HandlerFunc(app.logout)))
	mux.Get("/", app.session.Enable(app.authenticate(http.HandlerFunc(app.tasks))))
	mux.Post("/tasks", app.session.Enable(app.authenticate(http.HandlerFunc(app.tasks))))
	mux.Post("/delete", app.session.Enable(http.HandlerFunc(app.deleteTask)))
	mux.Post("/update", app.session.Enable(http.HandlerFunc(app.updateTask)))
	mux.Get("/panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Intentional panic for testing purposes")
	}))

	//struct for server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.LogRequest(app.secureHeaders((mux))),
	}

	//adding the static css file to fileserver
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	//listening to the port
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

// Function to Connect to DB
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
