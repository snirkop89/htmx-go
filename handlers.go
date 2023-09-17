package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/snirkop89/htmx/models"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (app *application) contacts(w http.ResponseWriter, r *http.Request) {
	var err error
	search := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("page")

	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	contacts, err := app.contactService.List(page, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Trigger") == "search" {
		render(w, r, "index.html", map[string]any{
			"Archiver": app.archiver,
			"Contacts": contacts,
			"Page":     page,
		}, "rows")
		return
	}
	render(w, r, "index.html", map[string]any{
		"Archiver": app.archiver,
		"Contacts": contacts,
		"Page":     page,
	})
}

func (app *application) newContact(w http.ResponseWriter, r *http.Request) {
	render(w, r, "new.html", map[string]any{
		"Contact": &models.Contact{},
		"Errors":  make(map[string]string),
	})
}

func (app *application) proccessNewContact(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	contact := models.Contact{
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		Email:     r.PostFormValue("email"),
		Phone:     r.PostFormValue("phone"),
	}
	_, err := app.contactService.Add(contact)
	if err != nil {
		render(w, r, "new.html", map[string]any{
			"Contact": contact,
			"Errors":  make(map[string]string),
		})
		return
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (app *application) viewContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	contact, err := app.contactService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	render(w, r, "show.html", map[string]any{
		"Contact": contact,
	})
}

func (app *application) validateEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	_, err := app.contactService.GetByEmail(email)
	if err == nil {
		w.Write([]byte("Email Must be unique"))
		return
	}

	w.Write([]byte(""))
}

func (app *application) editContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	contact, err := app.contactService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	render(w, r, "edit.html", map[string]any{
		"Contact": contact,
		"Errors":  make(map[string]string),
	})
}

func (app *application) processEditContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	r.ParseForm()

	contact := models.Contact{
		ID:        id,
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		Email:     r.PostFormValue("email"),
		Phone:     r.PostFormValue("phone"),
	}
	if err := app.contactService.Update(&contact); err != nil {
		render(w, r, "edit.html", map[string]any{
			"Contact": contact,
			"Errors":  make(map[string]string),
		})
		return
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (app *application) deleteContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	if err := app.contactService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Trigger") == "delete-btn" {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (app *application) contactsCount(w http.ResponseWriter, r *http.Request) {
	totalContacts := app.contactService.Count()
	fmt.Fprintf(w, " %d total counts", totalContacts)
}

func (app *application) deleteMultipleContacts(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	selectedIDs := r.PostForm["selected_contacts_ids"]
	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "unexpected error", http.StatusInternalServerError)
			return
		}
		err = app.contactService.Delete(id)
		if err != nil {
			http.Error(w, "unexpected error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (app *application) archiveContacts(w http.ResponseWriter, r *http.Request) {
	render(w, r, "index.html", map[string]any{
		"Archiver": app.archiver,
	}, "archive-ui")
}

func (app *application) deleteArchivedConctacts(w http.ResponseWriter, r *http.Request) {
	app.archiver.Reset()
	render(w, r, "index.html", map[string]any{
		"Archiver": app.archiver,
	}, "archive-ui")

}
func (app *application) processArchiveContacts(w http.ResponseWriter, r *http.Request) {
	app.archiver.Run()
	render(w, r, "index.html", map[string]any{
		"Archiver": app.archiver,
	}, "archive-ui")
}
