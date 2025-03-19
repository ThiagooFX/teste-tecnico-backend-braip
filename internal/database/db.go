package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Abrir a conexão com o banco SQLite
func OpenDB() (*sql.DB, error) {
	// Verifica se o arquivo do banco já existe, se não, cria ele
	_, err := os.Stat("database.db")
	if os.IsNotExist(err) {
		// Cria o banco de dados se não existir
		file, err := os.Create("database.db")
		if err != nil {
			log.Fatal("Não foi possível criar o banco de dados:", err)
			return nil, err
		}
		file.Close()
	}

	// Conectar ao banco de dados SQLite
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados:", err)
		return nil, err
	}

	// Verificar se a conexão foi bem-sucedida
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar a conexão:", err)
		return nil, err
	}

	DB = db
	return DB, nil
}

// Criar a tabela de produtos
func CreateTable() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price INTEGER NOT NULL,
		description TEXT NOT NULL,
		category TEXT NOT NULL,
		image_url TEXT
	);
	`
	_, err := DB.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Erro ao criar a tabela de produtos:", err)
	}
}
