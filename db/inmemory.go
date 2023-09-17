package db

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/snirkop89/htmx/models"
	"golang.org/x/exp/maps"
)

var ErrNotFound = errors.New("not found")

type InMemory struct {
	data map[int]*models.Contact
	mu   *sync.Mutex
}

func NewInMemory() *InMemory {
	data, err := os.ReadFile("db/contacts.json")
	if err != nil {
		panic(err)
	}
	var contacts []*models.Contact
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		panic(err)
	}

	db := make(map[int]*models.Contact)
	for i := range contacts {
		db[i] = contacts[i]
	}
	return &InMemory{
		data: db,
		mu:   &sync.Mutex{},
	}
}

func (d *InMemory) ByID(id int) (*models.Contact, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	c, ok := d.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (d *InMemory) ByEmail(email string) (*models.Contact, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, c := range d.data {
		if c.Email == email {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (d *InMemory) List(page int) ([]*models.Contact, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	start := (page - 1) * 10
	end := start + 10
	vals := maps.Values(d.data)
	if end > len(vals) {
		return vals[start:], nil
	}
	return vals[start:end], nil
}

func (d *InMemory) Search(searchTerm string) ([]*models.Contact, error) {
	var contacts []*models.Contact
	for _, contact := range d.data {
		if strings.Contains(strings.ToLower(contact.FirstName), searchTerm) ||
			strings.Contains(strings.ToLower(contact.LastName), searchTerm) ||
			strings.Contains(strings.ToLower(contact.Email), searchTerm) {
			contacts = append(contacts, contact)
		}
	}

	return contacts, nil
}

func (d *InMemory) Insert(contact models.Contact) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	newID := len(d.data) + 1
	contact.ID = newID

	d.data[newID] = &contact

	return newID, nil
}

func (d *InMemory) Update(contact *models.Contact) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.data[contact.ID] = contact
	return nil
}

func (d *InMemory) Delete(id int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.data, id)

	return nil
}

func (d *InMemory) Count() int {
	return len(d.data)
}
