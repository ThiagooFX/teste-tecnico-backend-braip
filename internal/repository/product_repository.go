package repository

import (
	"braip/internal/database"
	"braip/internal/models"
	"database/sql"
	"log"
)

// GetProducts retorna todos os produtos do banco de dados
func GetProducts() ([]models.Product, error) {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, price, description, category, image_url FROM products")
	if err != nil {
		log.Printf("Erro ao buscar produtos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.ImageURL); err != nil {
			log.Printf("Erro ao processar produto: %v", err)
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// CreateProduct insere um novo produto no banco de dados
func CreateProduct(product models.Product) error {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO products (name, price, description, category, image_url) VALUES (?, ?, ?, ?, ?)",
		product.Name, product.Price, product.Description, product.Category, product.ImageURL)
	if err != nil {
		log.Printf("Erro ao salvar produto: %v", err)
		return err
	}

	return nil
}

// GetProductByID retorna um produto pelo ID
func GetProductByID(id int) (*models.Product, error) {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		return nil, err
	}
	defer db.Close()

	var product models.Product
	row := db.QueryRow("SELECT id, name, price, description, category, image_url FROM products WHERE id = ?", id)
	err = row.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category, &product.ImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Produto não encontrado
		}
		log.Printf("Erro ao buscar produto: %v", err)
		return nil, err
	}

	return &product, nil
}

// UpdateProduct atualiza um produto no banco de dados
func UpdateProduct(id int, product models.Product) error {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE products SET name = ?, price = ?, description = ?, category = ?, image_url = ? WHERE id = ?",
		product.Name, product.Price, product.Description, product.Category, product.ImageURL, id)
	if err != nil {
		log.Printf("Erro ao atualizar produto: %v", err)
		return err
	}

	return nil
}

// DeleteProduct remove um produto do banco de dados
func DeleteProduct(id int) error {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		log.Printf("Erro ao excluir produto: %v", err)
		return err
	}

	return nil
}

// SearchProductsByNameAndCategory busca produtos por nome e categoria
func SearchProductsByNameAndCategory(name, category string) ([]models.Product, error) {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, name, price, description, category, image_url 
			  FROM products WHERE name LIKE ? AND category LIKE ?`
	rows, err := db.Query(query, "%"+name+"%", "%"+category+"%")
	if err != nil {
		log.Printf("Erro ao buscar produtos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.ImageURL); err != nil {
			log.Printf("Erro ao processar produto: %v", err)
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// SearchProductsByCategory busca produtos por categoria
func SearchProductsByCategory(category string) ([]models.Product, error) {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, name, price, description, category, image_url 
			  FROM products WHERE category LIKE ?`
	rows, err := db.Query(query, "%"+category+"%")
	if err != nil {
		log.Printf("Erro ao buscar produtos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.ImageURL); err != nil {
			log.Printf("Erro ao processar produto: %v", err)
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// SearchProductsByImage busca produtos com ou sem imagem
func SearchProductsByImage(hasImage bool) ([]models.Product, error) {
	db, err := db.OpenDB()
	if err != nil {
		log.Printf("Erro ao conectar ao banco de dados: %v", err)
		return nil, err
	}
	defer db.Close()

	var query string
	if hasImage {
		query = `SELECT id, name, price, description, category, image_url 
				  FROM products WHERE image_url IS NOT NULL AND image_url != ''`
	} else {
		query = `SELECT id, name, price, description, category, image_url 
				  FROM products WHERE image_url IS NULL OR image_url = ''`
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Erro ao buscar produtos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.ImageURL); err != nil {
			log.Printf("Erro ao processar produto: %v", err)
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}