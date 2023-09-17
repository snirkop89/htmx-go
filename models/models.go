package models

import "time"

type repository interface {
	ByID(int) (*Contact, error)
	ByEmail(string) (*Contact, error)
	List(int) ([]*Contact, error)
	Search(string) ([]*Contact, error)
	Count() int
	Insert(Contact) (int, error)
	Update(*Contact) error
	Delete(id int) error
}

type Contact struct {
	ID        int    `json:"id"`
	FirstName string `json:"first"`
	LastName  string `json:"last"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type ContactService struct {
	DB repository
}

func (s *ContactService) List(page int, searchTerm ...string) ([]*Contact, error) {
	var contacts []*Contact
	var err error
	if len(searchTerm) > 0 {
		contacts, err = s.DB.Search(searchTerm[0])
	} else {
		contacts, err = s.DB.List(page)
	}
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *ContactService) Add(contact Contact) (int, error) {
	id, err := s.DB.Insert(contact)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ContactService) GetByID(id int) (*Contact, error) {
	return s.DB.ByID(id)
}

func (s *ContactService) GetByEmail(email string) (*Contact, error) {
	return s.DB.ByEmail(email)
}

func (s *ContactService) Update(contact *Contact) error {
	return s.DB.Update(contact)
}

func (s *ContactService) Delete(id int) error {
	return s.DB.Delete(id)
}

func (s *ContactService) Count() int {
	// Add some delay to mimik heavy load to lazy load
	time.Sleep(1 * time.Second)
	return s.DB.Count()
}
