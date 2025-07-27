package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lovodia/Product.git/internal/db"
	"github.com/Lovodia/Product.git/internal/models"
	"github.com/gorilla/mux"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, price, category_id, created_at FROM products")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID, &p.CreatedAt)
		products = append(products, p)
	}
	json.NewEncoder(w).Encode(products)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var p models.Product
	err := db.DB.QueryRow("SELECT id, name, price, category_id, created_at FROM products WHERE id = $1",
		id).Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID, &p.CreatedAt)

	if err != nil {
		http.Error(w, "Product not found", 404)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	json.NewDecoder(r.Body).Decode(&p)

	err := db.DB.QueryRow(
		"INSERT INTO products(name, price, category_id, created_at) VALUES($1, $2, $3, NOW()) RETURNING id, created_at",
		p.Name, p.Price, p.CategoryID,
	).Scan(&p.ID, &p.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var p models.Product
	json.NewDecoder(r.Body).Decode(&p)

	_, err := db.DB.Exec(
		"UPDATE products SET name = $1, price = $2, category_id = $3 WHERE id = $4",
		p.Name, p.Price, p.CategoryID, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := db.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
