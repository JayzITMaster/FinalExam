package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"text/template"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	FromEmailAddress  string
	FromEmailPassword string
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data *templateData) {

	ts, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)

	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	// deal with the error status
	http.Error(w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionsManager.Exists(r.Context(), "authenticatedUserID")
}

func (app *application) isAuthenticatedAdmin(r *http.Request) bool {
	return app.sessionsManager.Exists(r.Context(), "authenticatedAdminID")
}
