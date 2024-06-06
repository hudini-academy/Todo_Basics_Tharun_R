package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (app *application) LogRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (app *application) authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !app.session.GetBool(r, "Authenticated") {
			app.session.Put(r, "flash", "Log In Before Accessing the resources")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (app *application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		fmt.Println(w.Header())

		rw := &responseWriter{
			ResponseWriter: w,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(rw, r)

		//fmt.Println(rw.body.String())
	})
}
