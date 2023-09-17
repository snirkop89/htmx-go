package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *application) http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/", app.homeHandler)
	r.Get("/contacts", app.contacts)
	r.Post("/contacts/delete", app.deleteMultipleContacts)
	r.Get("/contacts/count", app.contactsCount)
	r.Get("/contacts/archive", app.archiveContacts)
	r.Delete("/contacts/archive", app.deleteArchivedConctacts)
	r.Post("/contacts/archive", app.processArchiveContacts)
	r.Get("/contacts/{id}", app.viewContact)
	r.Get("/contacts/{id}/email", app.validateEmail)
	r.Get("/contacts/new", app.newContact)
	r.Post("/contacts/new", app.proccessNewContact)
	r.Get("/contacts/{id}/edit", app.editContact)
	r.Post("/contacts/{id}/edit", app.processEditContact)
	r.Delete("/contacts/{id}", app.deleteContact)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	return r
}
