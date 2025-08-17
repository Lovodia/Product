package httpDelivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/usecase"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	UC *usecase.CategoryUseCase
}

func NewCategoryHandler(uc *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{UC: uc}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.UC.GetAll()
	if err != nil {
		renderError(w, r, err)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	category, err := h.UC.GetByID(id)
	if err != nil {
		renderError(w, r, err)
		return
	}
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	id, err := h.UC.Create(category)
	if err != nil {
		renderError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	category.ID = id

	updated, err := h.UC.Update(category)
	if err != nil {
		renderError(w, r, err)
		return
	}
	if !updated {
		renderError(w, r, domain.ErrNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	deleted, err := h.UC.Delete(id)
	if err != nil {
		renderError(w, r, err)
		return
	}
	if !deleted {
		renderError(w, r, domain.ErrNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
