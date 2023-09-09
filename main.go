package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/", home)
	r.Get("/contacts", contacts)
	r.Get("/contacts/{id}", viewContact)
	r.Get("/contacts/new", newContact)
	r.Post("/contacts/new", proccessNewContact)
	r.Get("/contacts/{id}/edit", editContact)
	r.Post("/contacts/{id}/edit", processEditContact)
	r.Post("/contacts/{id}/delete", deleteContact)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func contacts(w http.ResponseWriter, r *http.Request) {
	search := chi.URLParam(r, "search")

	var contacts []*Contact
	var err error
	if search != "" {
		contacts, err = listContacts(search)
	} else {
		contacts, err = listContacts()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	render(w, r, "index.html", map[string]any{
		"Contacts": contacts,
	})
}

func newContact(w http.ResponseWriter, r *http.Request) {
	render(w, r, "new.html", map[string]any{
		"Contact": Contact{},
		"Errors":  make(map[string]string),
	})
}

func proccessNewContact(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	contact := Contact{
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		Email:     r.PostFormValue("email"),
		Phone:     r.PostFormValue("phone"),
	}
	if err := add(contact); err != nil {
		render(w, r, "new.html", map[string]any{
			"Contact": contact,
			"Errors":  make(map[string]string),
		})
		return
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func viewContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	contact, err := get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	render(w, r, "show.html", map[string]any{
		"Contact": contact,
	})
}

func editContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	contact, err := get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	render(w, r, "edit.html", map[string]any{
		"Contact": contact,
	})
}

func processEditContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	r.ParseForm()

	contact := Contact{
		ID:        id,
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		Email:     r.PostFormValue("email"),
		Phone:     r.PostFormValue("phone"),
	}
	if err := update(&contact); err != nil {
		render(w, r, "edit.html", map[string]any{
			"Contact": contact,
			"Errors":  make(map[string]string),
		})
		return
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	if err := delContact(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}
