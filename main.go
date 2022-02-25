package main

import (
	"github.com/joho/godotenv"
	"github.com/sgnl-05/contactService/api"
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

	http.HandleFunc("/api/list", h.ListContacts)
	http.HandleFunc("/api/add", h.AddContact)
	http.HandleFunc("/api/delete", h.DeleteContact)
	http.HandleFunc("/api/edit", h.EditContact)
	http.HandleFunc("/api/list-favs", h.ListFavorites)
	http.HandleFunc("/api/change-fav", h.ChangeFavorite)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
