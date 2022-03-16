package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/sgnl-05/contactService/api"
	"github.com/sgnl-05/contactService/storage"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var h api.ContactHandler
	api.ParseFlags(&h)

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
