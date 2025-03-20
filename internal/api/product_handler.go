package api

import (
	"encoding/json"
	"log"
	"net/http"
	"braip/internal/database"
	"github.com/gorilla/mux"
	"strconv"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       float64    `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
}


// Funçao para busca todos os produtos no banco de dados
func GetProducts(w http.ResponseWriter, r *http.Request) {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, price, description, category, image_url FROM products")
	if err != nil {
		log.Printf("Erro ao buscar produtos: %v", err)
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.ImageURL); err != nil {
			log.Printf("Erro ao processar produto: %v", err)
			http.Error(w, "Erro ao processar produto", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	// Verifica se não há produtos e retorna uma resposta adequada
	if len(products) == 0 {
		log.Println("Nenhum produto encontrado no banco de dados")
		http.Error(w, "Nenhum produto encontrado", http.StatusNotFound)
		return
	}

	// Retorna a lista de produtos como JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Printf("Erro ao converter produtos para JSON: %v", err)
		http.Error(w, "Erro ao processar resposta", http.StatusInternalServerError)
		return
	}
}

// Função para criar um novo produto
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validação dos dados
	if product.Name == "" || product.Price <= 0 || product.Description == "" || product.Category == "" {
		http.Error(w, "Campos obrigatórios não preenchidos corretamente", http.StatusBadRequest)
		return
	}

	// Salvar no banco de dados
	db, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}


	_, err = db.Exec("INSERT INTO products (name, price, description, category, image_url) VALUES (?, ?, ?, ?, ?)",
		product.Name, product.Price, product.Description, product.Category, product.ImageURL)
	if err != nil {
		http.Error(w, "Erro ao salvar produto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// Função para buscar um produto pelo ID

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("SELECT id, name, price, description, category, image_url FROM products WHERE id = ?", id)
	var product Product
	err = row.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category, &product.ImageURL)
	if err != nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// Função para atualizar um produto
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validação dos dados
	if product.Name == "" || product.Price <= 0 || product.Description == "" || product.Category == "" {
		http.Error(w, "Campos obrigatórios não preenchidos corretamente", http.StatusBadRequest)
		return
	}

	// Obter o ID do produto via URL
	params := mux.Vars(r)
	idStr := params["id"]

	// Converter ID para inteiro
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Atualizar no banco de dados
	db, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE products SET name = ?, price = ?, description = ?, category = ?, image_url = ? WHERE id = ?",
		product.Name, product.Price, product.Description, product.Category, product.ImageURL, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar produto", http.StatusInternalServerError)
		return
	}

	product.ID = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Função para deletar um produto
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Obter o ID do produto via URL
	params := mux.Vars(r)
	id := params["id"]

	// Deletar no banco de dados
	db, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Erro ao excluir produto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Função para buscar um produto de mesmo nome e categoria

func SearchProductsByNameAndCategory(w http.ResponseWriter, r *http.Request) {
	db, err := db.OpenDB()
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Recupera os parâmetros da URL
	name := r.URL.Query().Get("name")
	category := r.URL.Query().Get("category")

	if name == "" || category == "" {
		http.Error(w, "Parâmetros 'name' e 'category' são obrigatórios", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, price, description, category, image_url 
			  FROM products WHERE name LIKE ? AND category LIKE ?`

	rows, err := db.Query(query, "%"+name+"%", "%"+category+"%")
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.ImageURL); err != nil {
			http.Error(w, "Erro ao processar produto", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	// Retorna os produtos como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Função para buscar um produto de uma categoria em específico

func SearchProductsByCategory(w http.ResponseWriter, r *http.Request) {
	db, err := db.OpenDB()
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Recupera o parâmetro da URL
	category := r.URL.Query().Get("category")

	if category == "" {
		http.Error(w, "Parâmetro 'category' é obrigatório", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, price, description, category, image_url 
			  FROM products WHERE category LIKE ?`

	rows, err := db.Query(query, "%"+category+"%")
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.ImageURL); err != nil {
			http.Error(w, "Erro ao processar produto", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	// Retorna os produtos como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Função para buscar um produto se tiver imagem ou não 

func SearchProductsByImage(w http.ResponseWriter, r *http.Request) {
	db, err := db.OpenDB()
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Recupera o parâmetro da URL
	hasImage := r.URL.Query().Get("image")

	var query string
	if hasImage == "true" {
		query = `SELECT id, name, price, description, category, image_url 
				  FROM products WHERE image_url IS NOT NULL AND image_url != ''`
	} else if hasImage == "false" {
		query = `SELECT id, name, price, description, category, image_url 
				  FROM products WHERE image_url IS NULL OR image_url = ''`
	} else {
		http.Error(w, "Parâmetro 'image' deve ser 'true' ou 'false'", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.ImageURL); err != nil {
			http.Error(w, "Erro ao processar produto", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	// Retorna os produtos como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
