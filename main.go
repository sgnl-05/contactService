package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"sync"
)

var storageType *string

const filePath = "storage.json"

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

type contactHandler struct {
	mu      sync.Mutex
	storage StorageInterface
}

func (h *contactHandler) listContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendCustomError(w, http.StatusMethodNotAllowed, "GET requests only")
	}

	h.mu.Lock()
	allContacts, err := h.storage.List()
	h.mu.Unlock()

	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, fmt.Sprintf("Full list of contacts"), allContacts)
}

func (h *contactHandler) addContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendCustomError(w, http.StatusMethodNotAllowed, "POST requests only")
		return
	}
	if r.Header.Get("content-type") != "application/json" {
		sendCustomError(w, http.StatusUnsupportedMediaType, "JSON only")
		return
	}

	// Read request data
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close() //TODO Learn to catch errors in defer statements
	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var newContactBody Contact
	err = json.Unmarshal(body, &newContactBody)
	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Validation
	err = newContactBody.validate()
	if err != nil {
		sendCustomError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Filling missing values
	err = newContactBody.fillMissingFields()
	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Adding
	newContactBody.ID = uuid.New().String()
	err = h.storage.Add(newContactBody)
	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseBody := []Contact{newContactBody}
	sendSuccessResponse(w, fmt.Sprintf("New contact successfully added"), responseBody)
}

func (h *contactHandler) deleteContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendCustomError(w, http.StatusMethodNotAllowed, "GET requests only")
	}

	// Read request data
	keys := r.URL.Query()
	idDelete := keys.Get("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	// Deleting
	err := h.storage.Delete(idDelete)
	if err != nil {
		if errors.Is(err, ErrBReq) {
			sendCustomError(w, http.StatusBadRequest, fmt.Sprintf("No contact with ID: \"%v\"", idDelete))
			return
		}
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponseNoData(w, fmt.Sprintf("Contact \"%v\" successfully deleted", idDelete))
}

func (h *contactHandler) editContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendCustomError(w, http.StatusMethodNotAllowed, "POST requests only")
	}
	if r.Header.Get("content-type") != "application/json" {
		sendCustomError(w, http.StatusUnsupportedMediaType, "JSON only")
	}

	// Read request data
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var editContactBody EditContact
	err = json.Unmarshal(body, &editContactBody)
	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Validation
	err = editContactBody.validate()
	if err != nil {
		sendCustomError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Editing
	resultBody, err := h.storage.Edit(editContactBody)
	if err != nil {
		if errors.Is(err, ErrBReq) {
			sendCustomError(w, http.StatusBadRequest, fmt.Sprintf("No contact with ID: \"%v\"", editContactBody.ID))
			return
		}
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseBody := []Contact{resultBody}

	sendSuccessResponse(w, fmt.Sprintf("Contact \"%v\" successfully updated", editContactBody.ID), responseBody)
}

func (h *contactHandler) listFavorites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendCustomError(w, http.StatusMethodNotAllowed, "GET requests only")
	}

	h.mu.Lock()
	favContacts, err := h.storage.ListFavs()
	h.mu.Unlock()

	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, fmt.Sprintf("Full list of contacts"), favContacts)
}

func (h *contactHandler) changeFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendCustomError(w, http.StatusMethodNotAllowed, "GET requests only")
	}

	// Read params
	keys := r.URL.Query()
	id := keys.Get("id")
	if id == "" {
		sendCustomError(w, http.StatusBadRequest, "wrong request format, please use id={ID}&action={add|remove}")
		return
	}
	action := keys.Get("action")
	if action == "" {
		sendCustomError(w, http.StatusBadRequest, "wrong request format, please use id={ID}&action={add|remove}")
		return
	}

	//Changing
	h.mu.Lock()
	err := h.storage.ChangeFavs(id, action)
	h.mu.Unlock()
	if err != nil {
		if errors.Is(err, ErrAlreadyFav) {
			sendCustomError(w, http.StatusBadRequest, fmt.Sprintf("contact \"%v\" is already in favorites", id))
			return
		} else if errors.Is(err, ErrAlreadyNotFav) {
			sendCustomError(w, http.StatusBadRequest, fmt.Sprintf("contact \"%v\" is not in favorites already", id))
			return
		} else if errors.Is(err, ErrWrongFormat) {
			sendCustomError(w, http.StatusBadRequest, err.Error())
			return
		} else if errors.Is(err, ErrContactNotFound) {
			sendCustomError(w, http.StatusBadRequest, fmt.Sprintf("contact with ID \"%v\" not found", id))
			return
		} else {
			sendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if action == "add" {
		sendSuccessResponseNoData(w, fmt.Sprintf("Contact \"%v\" added to favorites", id))
		return
	}
	sendSuccessResponseNoData(w, fmt.Sprintf("Contact \"%v\" removed from favorites", id))
}

func main() {
	var h contactHandler

	/*h.storage = MemoryStorage{
		data: make(map[string]*Contact),
	}*/

	// h.storage = FileStorage{}

	parseFlags(&h)

	http.HandleFunc("/api/list", h.listContacts)
	http.HandleFunc("/api/add", h.addContact)
	http.HandleFunc("/api/delete", h.deleteContact)
	http.HandleFunc("/api/edit", h.editContact)
	http.HandleFunc("/api/list-favs", h.listFavorites)
	http.HandleFunc("/api/change-fav", h.changeFavorite)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
