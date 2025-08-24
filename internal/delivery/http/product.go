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

	r.HandleFunc("/products", errorHandler(handler.GetAll)).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}", errorHandler(handler.GetByID)).Methods(http.MethodGet)
	r.HandleFunc("/products", errorHandler(handler.Create)).Methods(http.MethodPost)
	r.HandleFunc("/products/{id}", errorHandler(handler.Update)).Methods(http.MethodPut)
	r.HandleFunc("/products/{id}", errorHandler(handler.Delete)).Methods(http.MethodDelete)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	products, err := h.UC.GetAll(r.Context())
	if err != nil {
		return err
	}
	renderJSON(w, http.StatusOK, products)
	return nil
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return domain.ErrBadRequest
	}

	product, err := h.UC.GetByID(r.Context(), id)
	if err != nil {
		return err
	}
	renderJSON(w, http.StatusOK, product)
	return nil
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return domain.ErrInvalidInput
	}

	id, err := h.UC.Create(r.Context(), product)
	if err != nil {
		return err
	}
	renderJSON(w, http.StatusCreated, map[string]int{"id": id})
	return nil
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return domain.ErrBadRequest
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return domain.ErrInvalidInput
	}

	updated, err := h.UC.Update(r.Context(), id, product)
	if err != nil {
		return err
	}
	if !updated {
		return domain.ErrNotFound
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) error {
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
