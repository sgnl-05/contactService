package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type StorageInterface interface {
	List() ([]Contact, error)
	Add(Contact) error
	Delete(string) error
	Edit(EditContact) (Contact, error)
	ListFavs() ([]Contact, error)
	ChangeFavs(string, string) error
}

type MemoryStorage struct {
	contactBook map[string]*Contact
}

type FileStorage struct{}

///////////////////////////////////////// MEMORY METHODS /////////////////////////////////////////

func (s MemoryStorage) List() ([]Contact, error) {
	var jsonContacts []Contact

	for _, v := range s.contactBook {
		jsonContacts = append(jsonContacts, *v)
	}

	return jsonContacts, nil
}

func (s MemoryStorage) Add(c Contact) error {
	s.contactBook[c.ID] = &c

	return nil
}

func (s MemoryStorage) Delete(id string) error {
	for k := range s.contactBook {
		if k == id {
			delete(s.contactBook, id)
			return nil
		}
	}

	return fmt.Errorf("no contact with ID: \"%v\"", id)
}

func (s MemoryStorage) Edit(e EditContact) (Contact, error) {
	var res Contact

	for k, v := range s.contactBook {
		if k == e.ID {
			if e.Name != "" {
				v.Name = e.Name
			}
			if e.Phone != "" {
				v.Phone = e.Phone
			}
			if e.Gender != "" {
				v.Gender = e.Gender
			}
			if e.Country != "" {
				v.Country = e.Country
			}
			res = *v
			return res, nil
		}
	}

	return res, fmt.Errorf("no contact with ID: \"%v\"", e.ID)
}

func (s MemoryStorage) ListFavs() ([]Contact, error) {
	var resultData []Contact

	for _, v := range s.contactBook {
		if v.Favorite == true {
			resultData = append(resultData, *v)
		}
	}

	return resultData, nil
}

func (s MemoryStorage) ChangeFavs(id string, action string) error {
	if _, ok := s.contactBook[id]; !ok {
		return ErrContactNotFound
	}

	switch action {
	case "add":
		if s.contactBook[id].Favorite == true {
			return ErrAlreadyFav
		}
		s.contactBook[id].Favorite = true
	case "remove":
		if s.contactBook[id].Favorite == false {
			return ErrAlreadyNotFav
		}
		s.contactBook[id].Favorite = false
	default:
		return ErrWrongFormat
	}

	return nil
}

//////////////////////////////////////// FILE METHODS ////////////////////////////////////////

func readFileContents() ([]Contact, error) {
	var contacts []Contact

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return contacts, err
	}

	err = json.Unmarshal(bytes, &contacts)
	return contacts, err
}

func writeFileContents(contacts []Contact) error {
	dataBytes, err := json.Marshal(contacts)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s FileStorage) List() ([]Contact, error) {
	contactList, err := readFileContents()

	return contactList, err
}

func (s FileStorage) Add(c Contact) error {
	contactList, err := readFileContents()
	if err != nil {
		return err
	}

	contactList = append(contactList, c)

	err = writeFileContents(contactList)
	if err != nil {
		return err
	}

	return nil
}

func (s FileStorage) Delete(id string) error {
	contactList, err := readFileContents()
	if err != nil {
		return err
	} // Internal

	for i := range contactList {
		if contactList[i].ID == id {
			contactList = append(contactList[:i], contactList[i+1:]...)
			err = writeFileContents(contactList)
			if err != nil {
				return err
			} // Internal
			return nil
		}
	}

	return ErrBReq // Bad request
}

func (s FileStorage) Edit(e EditContact) (Contact, error) {
	var res Contact

	contactList, err := readFileContents()
	if err != nil {
		return res, err
	} // Internal

	for i := range contactList {
		if e.ID == contactList[i].ID {
			if e.Name != "" {
				contactList[i].Name = e.Name
			}
			if e.Phone != "" {
				contactList[i].Phone = e.Phone
			}
			if e.Gender != "" {
				contactList[i].Gender = e.Gender
			}
			if e.Country != "" {
				contactList[i].Country = e.Country
			}
			res = contactList[i]
			break
		}
	}

	if res.ID == "" {
		return res, ErrBReq // Bad request
	}

	err = writeFileContents(contactList)
	if err != nil {
		return res, err
	} // Internal

	return res, nil
}

func (s FileStorage) ListFavs() ([]Contact, error) {
	var resultData []Contact
	contactList, err := readFileContents()

	for _, v := range contactList {
		if v.Favorite == true {
			resultData = append(resultData, v)
		}
	}

	return resultData, err
}

func (s FileStorage) ChangeFavs(id string, action string) error {
	contactList, err := readFileContents()
	if err != nil {
		return err
	} // Internal

	for i := range contactList {
		if contactList[i].ID == id {
			switch action {
			case "add":
				if contactList[i].Favorite == true {
					return ErrAlreadyFav
				}
				contactList[i].Favorite = true
				err = writeFileContents(contactList)
				if err != nil {
					return err
				} // Internal
				return nil
			case "remove":
				if contactList[i].Favorite == false {
					return ErrAlreadyNotFav
				}
				contactList[i].Favorite = false
				err = writeFileContents(contactList)
				if err != nil {
					return err
				} // Internal
				return nil
			default:
				return ErrWrongFormat
			}
		}
	}

	return ErrContactNotFound
}
