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

func NewProductHandler(r *mux.Router, uc *usecase.ProductUseCase) {
	handler := &ProductHandler{UC: uc}

	r.HandleFunc("/products", handler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}", handler.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/products", handler.Create).Methods(http.MethodPost)
	r.HandleFunc("/products/{id}", handler.Update).Methods(http.MethodPut)
	r.HandleFunc("/products/{id}", handler.Delete).Methods(http.MethodDelete)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.UC.GetAll()
	if err != nil {
		renderError(w, r, err)
		return
	}
	renderJSON(w, http.StatusOK, products)
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
	renderJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		renderError(w, r, domain.ErrInvalidImput)
		return
	}

	id, err := h.UC.Create(product)
	if err != nil {
		renderError(w, r, err)
		return
	}
	renderJSON(w, http.StatusCreated, map[string]int{"id": id})
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
		renderError(w, r, domain.ErrInvalidImput)
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
