package storage

import "github.com/sgnl-05/contactService/utils"

func (s MemoryStorage) List() ([]Contact, error) {
	var jsonContacts []Contact

	for _, v := range s.ContactBook {
		jsonContacts = append(jsonContacts, *v)
	}

	return jsonContacts, nil
}

func (s MemoryStorage) Add(c Contact) error {
	s.ContactBook[c.ID] = &c

	return nil
}

func (s MemoryStorage) Delete(id string) error {
	for k := range s.ContactBook {
		if k == id {
			delete(s.ContactBook, id)
			return nil
		}
	}

	return utils.ErrContactNotFound
}

func (s MemoryStorage) Edit(e EditContact) (Contact, error) {
	var res Contact

	for k, v := range s.ContactBook {
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

	return res, utils.ErrContactNotFound
}

func (s MemoryStorage) ListFavs() ([]Contact, error) {
	var resultData []Contact

	for _, v := range s.ContactBook {
		if v.Favorite == true {
			resultData = append(resultData, *v)
		}
	}

	return resultData, nil
}

func (s MemoryStorage) ChangeFavs(id string, action string) error {
	if _, ok := s.ContactBook[id]; !ok {
		return utils.ErrContactNotFound
	}

	switch action {
	case "add":
		if s.ContactBook[id].Favorite == true {
			return utils.ErrAlreadyFav
		}
		s.ContactBook[id].Favorite = true
	case "remove":
		if s.ContactBook[id].Favorite == false {
			return utils.ErrAlreadyNotFav
		}
		s.ContactBook[id].Favorite = false
	default:
		return utils.ErrWrongFormat
	}

	return nil
}
