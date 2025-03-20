package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"braip/internal/database"
)

// Estrutura do produto conforme a API externa
type ExternalProduct struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"` 		// Esse valor será convertido em centavos antes de ir para o database
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
}

var dbMutex sync.Mutex // Mutex para sincronizar o acesso ao banco de dados

func main() {
	// Definir a flag --id para importar um produto por ID
	idFlag := flag.Int("id", 0, "ID do produto a ser importado")
	flag.Parse()

	if *idFlag != 0 {
		// Se o ID for fornecido, importar um produto específico
		fmt.Printf("Importando produto de ID: %d...\n", *idFlag)
		err := ImportProductByID(*idFlag)
		if err != nil {
			log.Fatalf("Erro ao importar produto: %v", err)
		}
	} else {
		// Se não for fornecido ID, importar todos os produtos
		fmt.Println("Importando todos os produtos...")
		err := ImportProducts()
		if err != nil {
			log.Fatalf("Erro ao importar produtos: %v", err)
		}
	}

	fmt.Println("Importação concluída com sucesso!")
}

// ImportProducts busca todos os produtos da API externa e insere no banco
func ImportProducts() error {
	resp, err := http.Get("https://fakestoreapi.com/products")
	if err != nil {
		return fmt.Errorf("erro ao buscar produtos da API externa: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	var apiProducts []ExternalProduct
	err = json.NewDecoder(resp.Body).Decode(&apiProducts)
	if err != nil {
		return fmt.Errorf("erro ao decodificar resposta da API: %v", err)
	}

	// Conectar ao banco de dados
	db, err := db.OpenDB()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	var wg sync.WaitGroup

	// Inserir cada produto no banco de dados em paralelo usando goroutines
	for _, apiProduct := range apiProducts {
		wg.Add(1)
		go func(apiProduct ExternalProduct) {
			defer wg.Done()
			// Usar o mutex para garantir que apenas uma goroutine acesse o banco de dados por vez
			dbMutex.Lock()
			defer dbMutex.Unlock()

			// Converter o preço de reais para centavos
			priceInCents := int(apiProduct.Price * 100)

			_, err := db.Exec(`
				INSERT INTO products (id, name, price, description, category, image_url)
				VALUES (?, ?, ?, ?, ?, ?)
				ON CONFLICT(id) DO NOTHING;`, // Usa ON CONFLICT para evitar duplicação
				apiProduct.ID, apiProduct.Title, priceInCents, apiProduct.Description, apiProduct.Category, apiProduct.Image,
			)
			if err != nil {
				log.Printf("Erro ao inserir produto %d: %v", apiProduct.ID, err)
			} else {
				fmt.Printf("Produto %d importado com sucesso!\n", apiProduct.ID)
			}
		}(apiProduct)
	}

	// Aguardar até todas as goroutines terminarem
	wg.Wait()

	fmt.Println("Produtos importados com sucesso!")
	return nil
}

// ImportProductByID busca um produto específico da API externa e insere no banco
func ImportProductByID(id int) error {
	// Buscar produto pela API externa
	url := fmt.Sprintf("https://fakestoreapi.com/products/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("erro ao buscar produto pela API externa: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	var apiProduct ExternalProduct
	err = json.NewDecoder(resp.Body).Decode(&apiProduct)
	if err != nil {
		return fmt.Errorf("erro ao decodificar resposta da API: %v", err)
	}

	// Conectar ao banco de dados
	db, err := db.OpenDB()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Converter o preço de reais para centavos
	priceInCents := int(apiProduct.Price * 100)

	// Inserir o produto no banco
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err = db.Exec(`
		INSERT INTO products (id, name, price, description, category, image_url)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO NOTHING;`, // Usa ON CONFLICT para evitar duplicação
		apiProduct.ID, apiProduct.Title, priceInCents, apiProduct.Description, apiProduct.Category, apiProduct.Image,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir produto %d no banco de dados: %v", apiProduct.ID, err)
	}

	fmt.Printf("Produto %d importado com sucesso!\n", apiProduct.ID)
	return nil
}