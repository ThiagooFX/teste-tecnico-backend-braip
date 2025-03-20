package main

import (
	"log"
	"braip/internal/database"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"braip/internal/api"
)

func main() {
	// Abrir conexão com o banco
	_, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	// Criar a tabela se necessário
	db.CreateTable()
	log.Println("Banco de dados e tabela criados com sucesso!")


	// Rotas da API

	r := mux.NewRouter()

	// Rotas de consulta de produtos
	r.HandleFunc("/products/{id}", api.GetProductByID).Methods("GET")											// OK
	r.HandleFunc("/products/search/categoryandname", api.SearchProductsByNameAndCategory).Methods("GET")		// OK
	r.HandleFunc("/products/search/category", api.SearchProductsByCategory).Methods("GET")						// OK
	r.HandleFunc("/products/search/image", api.SearchProductsByImage).Methods("GET")							// OK

	//Demais rotas padrões do CRUD
	r.HandleFunc("/products", api.CreateProduct).Methods("POST")								// OK
	r.HandleFunc("/products", api.GetProducts).Methods("GET")									// OK
	r.HandleFunc("/products/{id}", api.UpdateProduct).Methods("PUT")							// OK
	r.HandleFunc("/products/{id}", api.DeleteProduct).Methods("DELETE")							// OK

	fmt.Println("Servidor rodando na porta 4000...")
	log.Fatal(http.ListenAndServe(":4000", r))


}