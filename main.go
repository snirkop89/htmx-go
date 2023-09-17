package main

import (
	"log"
	"net/http"

	"github.com/snirkop89/htmx/archiver"
	"github.com/snirkop89/htmx/db"
	"github.com/snirkop89/htmx/models"
)

type application struct {
	contactService *models.ContactService
	archiver       *archiver.Archiver
}

func main() {
	app := application{
		contactService: &models.ContactService{
			DB: db.NewInMemory(),
		},
		archiver: archiver.New(),
	}
	if err := http.ListenAndServe(":8080", routes(&app)); err != nil {
		log.Fatal(err)
	}
}
