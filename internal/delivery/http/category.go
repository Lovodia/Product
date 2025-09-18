package httpdelivery

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

func NewCategoryHandler(r *mux.Router, uc *usecase.CategoryUseCase) {
	handler := &CategoryHandler{UC: uc}

	r.HandleFunc("/categories", errorHandler(handler.GetAll)).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}", errorHandler(handler.GetByID)).Methods(http.MethodGet)
	r.HandleFunc("/categories", errorHandler(handler.Create)).Methods(http.MethodPost)
	r.HandleFunc("/categories/{id}", errorHandler(handler.Update)).Methods(http.MethodPut)
	r.HandleFunc("/categories/{id}", errorHandler(handler.Delete)).Methods(http.MethodDelete)
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	categories, err := h.UC.GetAll(r.Context())
	if err != nil {
		return err
	}

	renderJSON(w, http.StatusOK, categories)

	return nil
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return domain.ErrBadRequest
	}

	category, err := h.UC.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	renderJSON(w, http.StatusOK, category)

	return nil
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		return domain.ErrInvalidInput
	}

	id, err := h.UC.Create(r.Context(), category)
	if err != nil {
		return err
	}

	renderJSON(w, http.StatusCreated, map[string]int{"id": id})

	return nil
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return domain.ErrBadRequest
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		return domain.ErrInvalidInput
	}

	category.ID = id

	updated, err := h.UC.Update(r.Context(), category)
	if err != nil {
		return err
	}

	if !updated {
		return domain.ErrNotFound
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return domain.ErrBadRequest
	}

	deleted, err := h.UC.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	if !deleted {
		return domain.ErrNotFound
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
