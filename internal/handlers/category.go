package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lovodia/Product.git/internal/db"
	"github.com/Lovodia/Product.git/internal/models"
	"github.com/gorilla/mux"
)

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name)
		categories = append(categories, c)
	}

	json.NewEncoder(w).Encode(categories)
}

func GetCategoriesByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var c models.Category
	err := db.DB.QueryRow("SELECT id, name FROM categories WHIRE id = $1", id).Scan(&c.ID, &c.Name)
	if err != nil {
		http.Error(w, "category not found", 404)
		return
	}
	json.NewEncoder(w).Encode(c)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	json.NewDecoder(r.Body).Decode(&c)

	err := db.DB.QueryRow("INSERT INTO categories(name) VALUES($1) RETURNING id", c.Name).Scan(&c.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(c)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var c models.Category
	json.NewDecoder(r.Body).Decode(&c)

	_, err := db.DB.Exec("UPDATE categories SET name=$1 WHERE id=$2", c.Name, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := db.DB.Exec("DELETE FROM categories WHETE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
