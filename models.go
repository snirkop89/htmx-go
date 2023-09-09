package main

import "fmt"

type Contact struct {
	ID        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

var db = []*Contact{
	{
		ID:        1,
		FirstName: "John",
		LastName:  "Cohn",
		Phone:     "555-555555",
		Email:     "john@there.com",
	},
	{
		ID:        2,
		FirstName: "Dana",
		LastName:  "Crandith",
		Phone:     "12-456-789",
		Email:     "drcan@example.com",
	},
}

func listContacts(searchTerm ...string) ([]*Contact, error) {
	return db, nil
}

func add(contact Contact) error {
	if contact.ID == 0 {
		contact.ID = len(db) + 1
	}
	db = append(db, &contact)
	return nil
}

func get(id int) (*Contact, error) {
	for _, c := range db {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("contact not found")
}

func update(contact *Contact) error {
	for _, c := range db {
		if c.ID == contact.ID {
			c.Email = contact.Email
			c.FirstName = contact.FirstName
			c.LastName = contact.LastName
			c.Phone = contact.Phone
		}
	}
	return fmt.Errorf("contact not found")
}

func delContact(id int) error {
	for i := range db {
		if db[i].ID == id {
			db = append(db[:i], db[i+1:]...)
			break
		}
	}
	return nil
}
