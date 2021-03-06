package storage

import (
	"encoding/json"
	"github.com/sgnl-05/contactService/utils"
	"io/ioutil"
	"os"
	"strings"
)

func readFileContents() ([]Contact, error) {
	var contacts []Contact
	filePath := os.Getenv("LOCAL_FILENAME")

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return contacts, err
	}

	if len(bytes) == 0 {
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

	filePath := os.Getenv("LOCAL_FILENAME")
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

	return utils.ErrContactNotFound // Bad request
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
		return res, utils.ErrContactNotFound // Bad request
	}

	err = writeFileContents(contactList)
	if err != nil {
		return res, err
	} // Internal

	return res, nil
}

func (s FileStorage) Filter(field string, value string) ([]Contact, error) {
	var resultData []Contact
	fullList, err := readFileContents()
	if err != nil {
		return resultData, err
	}

	switch field {
	case "name":
		for _, v := range fullList {
			if strings.Contains(
				strings.ToLower(v.Name),
				strings.ToLower(value),
			) {
				resultData = append(resultData, v)
			}
		}
		return resultData, nil
	case "phone":
		for _, v := range fullList {
			if strings.Contains(
				strings.ToLower(v.Phone),
				strings.ToLower(value),
			) {
				resultData = append(resultData, v)
			}
		}
		return resultData, nil
	default:
		return resultData, utils.ErrFilterWrongFormat
	}
}

func (s FileStorage) ListFavs() ([]Contact, error) {
	var resultData []Contact
	contactList, err := readFileContents()
	if err != nil {
		return resultData, err
	}

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
					return utils.ErrAlreadyFav
				}
				contactList[i].Favorite = true
				err = writeFileContents(contactList)
				if err != nil {
					return err
				} // Internal
				return nil
			case "remove":
				if contactList[i].Favorite == false {
					return utils.ErrAlreadyNotFav
				}
				contactList[i].Favorite = false
				err = writeFileContents(contactList)
				if err != nil {
					return err
				} // Internal
				return nil
			default:
				return utils.ErrFavWrongFormat
			}
		}
	}

	return utils.ErrContactNotFound
}
