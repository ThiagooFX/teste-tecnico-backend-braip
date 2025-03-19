package main

import (
	"log"
	"braip/internal/database"
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
}