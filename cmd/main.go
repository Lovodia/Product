package main

import (
	"log"
	"net/http"

	"github.com/Lovodia/Product.git/internal/db"
	"github.com/Lovodia/Product.git/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatal("Error connecting to DB", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/categories", handlers.GetAllCategories).Methods("GET")
	r.HandleFunc("/categories/{id:[0-9]+}", handlers.GetCategoriesByID).Methods("GET")
	r.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")
	r.HandleFunc("/categories/{id:[0-9]+}", handlers.UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id:[0-9]+}", handlers.DeleteCategory).Methods("DELETE")

	r.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", handlers.GetProductByID).Methods("GET")
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id:[0-9]+}", handlers.DeleteProduct).Methods("DELETE")

	log.Println("The server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
