package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"github.com/sgnl-05/contactService/api"
	"github.com/sgnl-05/contactService/storage"
	"log"
	"net/http"
	"os"
)

type options struct {
	StorageType string `short:"d" description:"Data storage type" choice:"memory" choice:"file" choice:"elastic" required:"true"`
}

func parseFlags(h *api.ContactHandler, o options) {
	switch o.StorageType {
	case "memory":
		fmt.Println("Store in memory")
		h.Storage = storage.MemoryStorage{
			ContactBook: make(map[string]*storage.Contact),
		}
	case "file":
		fmt.Println("Store in local file")
		h.Storage = storage.FileStorage{}
	case "elastic":
		fmt.Println("Store in Elastic")
		h.Storage = storage.NewElasticStorage()
	default:
		fmt.Println("Available -d key values: memory|file|elastic")
		return
	}
}

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var h api.ContactHandler
	var o options

	if _, err := flags.Parse(&o); err != nil {
		fmt.Printf("Parse args error %+v", err)
		os.Exit(1)
	}

	parseFlags(&h, o)

	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.SetHeader("content-type", "application/json"))

	r.Route("/api", func(r chi.Router) {
		r.Get("/list", h.ListContacts)
		r.Get("/delete", h.DeleteContact)
		r.With(storage.ValidateNewContact).Post("/add", h.AddContact)
		r.With(storage.ValidateExistingContact).Post("/edit", h.EditContact)
		r.Post("/filter", h.Filter)
		r.Get("/list-favs", h.ListFavorites)
		r.Get("/change-fav", h.ChangeFavorite)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
