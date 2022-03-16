package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sgnl-05/contactService/utils"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Genderize struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type CoProb struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Nationalize struct {
	Name    string   `json:"name"`
	Country []CoProb `json:"country"`
}

func (c *Contact) FillMissingFields() error {
	var err error

	if c.Gender == "" {
		err = c.genderize()
		if err != nil {
			return err
		}
	}

	if c.Country == "" {
		err = c.nationalize()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Contact) genderize() error {
	nameUrl := fmt.Sprintf("https://api.genderize.io?name=%v", c.Name)
	resp, err := http.Get(nameUrl)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	var gResponseBody Genderize
	err = json.Unmarshal(body, &gResponseBody)
	if err != nil {
		return err
	}

	c.Gender = gResponseBody.Gender

	return nil
}

func (c *Contact) nationalize() error {
	nameUrl := fmt.Sprintf("https://api.nationalize.io?name=%v", c.Name)
	resp, err := http.Get(nameUrl)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	var nResponseBody Nationalize
	err = json.Unmarshal(body, &nResponseBody)
	if err != nil {
		return err
	}

	/*
		{"name":"michael", "country":[
		{"country_id":"US","probability":0.08986482266532715},
		{"country_id":"AU","probability":0.05976757527083082},
		{"country_id":"NZ","probability":0.04666974820852911}
		]
		}
	*/

	highProb := 0.0
	resCountry := ""
	for _, v := range nResponseBody.Country {
		if v.Probability > highProb {
			highProb = v.Probability
			resCountry = v.CountryID
		}
	}

	c.Country = resCountry

	return nil
}

func validateName(name string) error {
	if name == "" || len(name) < 4 {
		return errors.New("name must have more than 4 characters")
	}

	return nil
}

func validatePhone(phone string) error {
	phoneReg, err := regexp.MatchString(`^\+7\d{10}$`, phone)
	if phoneReg == false || err != nil {
		return errors.New("phone number must be in \"+7xxxxxxxxxx\" format")
	}

	return nil
}

func validateGender(gender string) error {
	if gender != "male" && gender != "female" && gender != "" {
		return errors.New("gender must be either \"male\" or \"female\", liberal")
	}

	return nil
}

func validateCountry(country string) error {
	if country == "" {
		return nil
	}

	countryReg, err := regexp.MatchString(`^[A-Z]{2}`, country)
	if countryReg == false || err != nil {
		return errors.New("country code must consist of two uppercase letters")
	}

	return nil
}

func ValidateNewContact(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var c Contact
		err = json.Unmarshal(body, &c)
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = validateName(c.Name)
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = validatePhone(c.Phone)
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = validateGender(c.Gender)
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = validateCountry(c.Country)
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}

		//err = r.Body.Close()
		err = r.Body.Close()
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(w, r)
	})
}

func ValidateExistingContact(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var c Contact
		err = json.Unmarshal(body, &c)
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if c.Name != "" {
			err = validateName(c.Name)
			if err != nil {
				utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if c.Phone != "" {
			err = validatePhone(c.Phone)
			if err != nil {
				utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if c.Gender != "" {
			err = validateGender(c.Gender)
			if err != nil {
				utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if c.Country != "" {
			err = validateCountry(c.Country)
			if err != nil {
				utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		//err = r.Body.Close()
		err = r.Body.Close()
		if err != nil {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(w, r)
	})
}

func (c *EditContact) Validate() error {
	var err error

	if c.Name != "" {
		err = validateName(c.Name)
		if err != nil {
			return err
		}
	}

	if c.Phone != "" {
		err = validatePhone(c.Phone)
		if err != nil {
			return err
		}
	}

	if c.Gender != "" {
		err = validateGender(c.Gender)
		if err != nil {
			return err
		}
	}

	if c.Country != "" {
		err = validateCountry(c.Country)
		if err != nil {
			return err
		}
	}

	return nil
}
