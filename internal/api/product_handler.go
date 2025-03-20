package api

import (
	"braip/internal/models"
	"braip/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetProducts retorna todos os produtos
func GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := services.GetProducts()
	if err != nil {
		log.Printf("Erro ao buscar produtos: %v", err)
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	if len(products) == 0 {
		log.Println("Nenhum produto encontrado")
		http.Error(w, "Nenhum produto encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// CreateProduct cria um novo produto
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if product.Name == "" || product.Price <= 0 || product.Description == "" || product.Category == "" {
		http.Error(w, "Campos obrigatórios não preenchidos corretamente", http.StatusBadRequest)
		return
	}

	if err := services.CreateProduct(product); err != nil {
		http.Error(w, "Erro ao criar produto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetProductByID retorna um produto pelo ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	product, err := services.GetProductByID(id)
	if err != nil {
		http.Error(w, "Erro ao buscar produto", http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct atualiza um produto
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if product.Name == "" || product.Price <= 0 || product.Description == "" || product.Category == "" {
		http.Error(w, "Campos obrigatórios não preenchidos corretamente", http.StatusBadRequest)
		return
	}

	if err := services.UpdateProduct(id, product); err != nil {
		http.Error(w, "Erro ao atualizar produto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct remove um produto
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := services.DeleteProduct(id); err != nil {
		http.Error(w, "Erro ao excluir produto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchProductsByNameAndCategory busca produtos por nome e categoria
func SearchProductsByNameAndCategory(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	category := r.URL.Query().Get("category")

	if name == "" || category == "" {
		http.Error(w, "Parâmetros 'name' e 'category' são obrigatórios", http.StatusBadRequest)
		return
	}

	products, err := services.SearchProductsByNameAndCategory(name, category)
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// SearchProductsByCategory busca produtos por categoria
func SearchProductsByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	if category == "" {
		http.Error(w, "Parâmetro 'category' é obrigatório", http.StatusBadRequest)
		return
	}

	products, err := services.SearchProductsByCategory(category)
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// SearchProductsByImage busca produtos com ou sem imagem
func SearchProductsByImage(w http.ResponseWriter, r *http.Request) {
	hasImage := r.URL.Query().Get("image")

	var hasImageBool bool
	if hasImage == "true" {
		hasImageBool = true
	} else if hasImage == "false" {
		hasImageBool = false
	} else {
		http.Error(w, "Parâmetro 'image' deve ser 'true' ou 'false'", http.StatusBadRequest)
		return
	}

	products, err := services.SearchProductsByImage(hasImageBool)
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}