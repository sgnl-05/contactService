package storage

import (
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Contact struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Country  string `json:"country"`
	Favorite bool   `json:"favorite"`
}

type EditContact struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Gender  string `json:"gender"`
	Country string `json:"country"`
}

type FilterRequest struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type StorageInterface interface {
	List() ([]Contact, error)
	Add(Contact) error
	Delete(string) error
	Edit(EditContact) (Contact, error)
	Filter(string, string) ([]Contact, error)
	ListFavs() ([]Contact, error)
	ChangeFavs(string, string) error
}

type MemoryStorage struct {
	ContactBook map[string]*Contact
}

type FileStorage struct{}

const IndexName = "contacts"

type ElasticStorage struct {
	client *elasticsearch.Client
}

func NewElasticStorage() ElasticStorage {
	var esObject ElasticStorage

	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTIC_URL"),
		},
		Username: os.Getenv("ELASTIC_USERNAME"),
		Password: os.Getenv("ELASTIC_PASSWORD"),
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // TODO Probably shouldn't do this
				MinVersion:         tls.VersionTLS11,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	esObject.client = es
	return esObject
}
