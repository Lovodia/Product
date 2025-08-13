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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	category, err := h.UC.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid imput", http.StatusBadRequest)
		return
	}

	id, err := h.UC.Create(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	category.ID = id

	updated, err := h.UC.Update(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !updated {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	deleted, err := h.UC.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
