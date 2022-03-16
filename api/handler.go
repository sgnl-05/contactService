package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"sync"

	"github.com/sgnl-05/contactService/storage"
	"github.com/sgnl-05/contactService/utils"
)

type ContactHandler struct {
	mu      sync.Mutex
	Storage storage.StorageInterface
}

func (h *ContactHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	allContacts, err := h.Storage.List()

	h.mu.Unlock()

	if err != nil {
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, fmt.Sprintf("Full list of contacts"), allContacts)
}

func (h *ContactHandler) AddContact(w http.ResponseWriter, r *http.Request) {
	// Read request data
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close() //TODO Learn to catch errors in defer statements
	if err != nil {
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var newContactBody storage.Contact
	err = json.Unmarshal(body, &newContactBody)
	if err != nil {
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Filling missing values
	err = newContactBody.FillMissingFields()
	if err != nil {
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Adding
	newContactBody.ID = uuid.New().String()
	err = h.Storage.Add(newContactBody)
	if err != nil {
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseBody := []storage.Contact{newContactBody}
	utils.SendSuccessResponse(w, fmt.Sprintf("New contact successfully added"), responseBody)
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	// Read request data
	keys := r.URL.Query()
	idDelete := keys.Get("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	// Deleting
	err := h.Storage.Delete(idDelete)
	if err != nil {
		if errors.Is(err, utils.ErrContactNotFound) {
			utils.SendCustomError(w, http.StatusBadRequest, fmt.Sprintf("No contact with ID: \"%v\"", idDelete))
			return
		}
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponseNoData(w, fmt.Sprintf("Contact \"%v\" successfully deleted", idDelete))
}

func (h *ContactHandler) EditContact(w http.ResponseWriter, r *http.Request) {
	// Read request data
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var editContactBody storage.EditContact
	err = json.Unmarshal(body, &editContactBody)
	if err != nil {
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Editing
	resultBody, err := h.Storage.Edit(editContactBody)
	if err != nil {
		if errors.Is(err, utils.ErrContactNotFound) {
			utils.SendCustomError(w, http.StatusBadRequest, fmt.Sprintf("No contact with ID: \"%v\"", editContactBody.ID))
			return
		}
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseBody := []storage.Contact{resultBody}

	utils.SendSuccessResponse(w, fmt.Sprintf("Contact \"%v\" successfully updated", editContactBody.ID), responseBody)
}

func (h *ContactHandler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	favContacts, err := h.Storage.ListFavs()
	h.mu.Unlock()

	if err != nil {
		utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, fmt.Sprintf("Full list of favorites"), favContacts)
}

func (h *ContactHandler) ChangeFavorite(w http.ResponseWriter, r *http.Request) {
	// Read params
	keys := r.URL.Query()
	id := keys.Get("id")
	if id == "" {
		utils.SendCustomError(w, http.StatusBadRequest, "wrong request format, please use id={ID}&action={add|remove}")
		return
	}
	action := keys.Get("action")
	if action == "" {
		utils.SendCustomError(w, http.StatusBadRequest, "wrong request format, please use id={ID}&action={add|remove}")
		return
	}

	//Changing
	h.mu.Lock()
	err := h.Storage.ChangeFavs(id, action)
	h.mu.Unlock()
	if err != nil {
		if errors.Is(err, utils.ErrAlreadyFav) {
			utils.SendCustomError(w, http.StatusBadRequest, fmt.Sprintf("contact \"%v\" is already in favorites", id))
			return
		} else if errors.Is(err, utils.ErrAlreadyNotFav) {
			utils.SendCustomError(w, http.StatusBadRequest, fmt.Sprintf("contact \"%v\" is not in favorites already", id))
			return
		} else if errors.Is(err, utils.ErrWrongFormat) {
			utils.SendCustomError(w, http.StatusBadRequest, err.Error())
			return
		} else if errors.Is(err, utils.ErrContactNotFound) {
			utils.SendCustomError(w, http.StatusBadRequest, fmt.Sprintf("contact with ID \"%v\" not found", id))
			return
		} else {
			utils.SendCustomError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if action == "add" {
		utils.SendSuccessResponseNoData(w, fmt.Sprintf("Contact \"%v\" added to favorites", id))
		return
	}
	utils.SendSuccessResponseNoData(w, fmt.Sprintf("Contact \"%v\" removed from favorites", id))
}
