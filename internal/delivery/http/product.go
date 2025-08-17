package httpDelivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/usecase"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	UC *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{UC: uc}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.UC.GetAll()
	if err != nil {
		renderError(w, r, err)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	product, err := h.UC.GetByID(id)
	if err != nil {
		renderError(w, r, err)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	id, err := h.UC.Create(product)
	if err != nil {
		renderError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		renderError(w, r, domain.ErrBadRequest)
		return
	}

	updated, err := h.UC.Update(id, product)
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

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
