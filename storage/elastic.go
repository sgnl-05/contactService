package storage

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/sgnl-05/contactService/utils"
	"strings"
)

type eContactSource struct {
	Source Contact `json:"_source"`
}

type eContactHitsLow struct {
	Hits []eContactSource `json:"hits"`
}

type eContactHitsUp struct {
	Hits eContactHitsLow `json:"hits"`
}

func (s ElasticStorage) updateElasticDoc(body Contact) error {
	contactString, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = s.client.Index(IndexName, strings.NewReader(string(contactString)), s.client.Index.WithDocumentID(body.ID))
	if err != nil {
		return err
	}

	return nil
}

func (s ElasticStorage) List() ([]Contact, error) {
	var list []Contact

	response, err := s.client.Search(
		s.client.Search.WithIndex(IndexName),
		s.client.Search.WithBody(strings.NewReader(`{
	  "query": {
	    "match_all": {}
	  }
	}`)),
	)
	if err != nil {
		return list, err
	}

	var responseBody eContactHitsUp
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return list, err
	}

	for i := range responseBody.Hits.Hits {
		list = append(list, responseBody.Hits.Hits[i].Source)
	}

	return list, nil
}

func (s ElasticStorage) Add(c Contact) error {
	contactString, err := json.Marshal(c)
	if err != nil {
		return err
	}

	request := esapi.IndexRequest{Index: IndexName, DocumentID: c.ID, OpType: "create", Body: strings.NewReader(string(contactString))}
	_, err = request.Do(context.Background(), s.client)
	if err != nil {
		return err
	}

	return nil
}

func (s ElasticStorage) Delete(id string) error {
	_, err := s.client.Delete(IndexName, id)

	if err != nil {
		return err
	}

	return nil
}

func (s ElasticStorage) Edit(e EditContact) (Contact, error) {
	var res Contact
	var responseBody eContactSource

	response, err := s.client.Get(IndexName, e.ID)
	if err != nil {
		return res, err
	}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return res, err
	}

	res = responseBody.Source

	if e.Name != "" {
		res.Name = e.Name
	}
	if e.Country != "" {
		res.Country = e.Country
	}
	if e.Phone != "" {
		res.Phone = e.Phone
	}
	if e.Gender != "" {
		res.Gender = e.Gender
	}

	err = s.updateElasticDoc(res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s ElasticStorage) Filter(field string, value string) ([]Contact, error) {
	var res []Contact
	var searchCond string

	switch field {
	case "name":
		searchCond = `{
	"query": {
		"regexp": {
			"name": ".*` + value + `.*"
			}
		}
	}`
	case "phone":
		searchCond = `{
	"query": {
		"regexp": {
			"phone": ".*` + value + `.*"
			}
		}
	}`
	default:
		return res, utils.ErrFilterWrongFormat
	}

	response, err := s.client.Search(
		s.client.Search.WithIndex(IndexName),
		s.client.Search.WithBody(strings.NewReader(searchCond)),
	)
	if err != nil {
		return res, err
	}

	var responseBody eContactHitsUp
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return res, err
	}

	for i := range responseBody.Hits.Hits {
		res = append(res, responseBody.Hits.Hits[i].Source)
	}

	return res, nil
}

func (s ElasticStorage) ListFavs() ([]Contact, error) {
	var list []Contact

	response, err := s.client.Search(
		s.client.Search.WithIndex(IndexName),
		s.client.Search.WithBody(strings.NewReader(`{
	  "query": {
	    "match": {
			"favorite": true
		}
	  }
	}`)),
	)
	if err != nil {
		return list, err
	}

	var responseBody eContactHitsUp
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return list, err
	}

	for i := range responseBody.Hits.Hits {
		list = append(list, responseBody.Hits.Hits[i].Source)
	}

	return list, nil
}

func (s ElasticStorage) ChangeFavs(id string, action string) error {
	var res Contact
	var responseBody eContactSource

	response, err := s.client.Get(IndexName, id)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	res = responseBody.Source
	if res.ID == id {
		switch action {
		case "add":
			if res.Favorite == true {
				return utils.ErrAlreadyFav
			}
			res.Favorite = true
			err = s.updateElasticDoc(res)
			if err != nil {
				return err
			} // Internal
			return nil
		case "remove":
			if res.Favorite == false {
				return utils.ErrAlreadyNotFav
			}
			res.Favorite = false
			err = s.updateElasticDoc(res)
			if err != nil {
				return err
			} // Internal
			return nil
		default:
			return utils.ErrFavWrongFormat
		}
	}

	return utils.ErrContactNotFound
}
