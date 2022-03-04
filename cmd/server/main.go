package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eoliveira-daitan/go-challenges/internal/api"
	"github.com/eoliveira-daitan/go-challenges/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Create a .env file in the root directory. Use `.env.example` as a start point")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	var repo repository.Repository

	cfg := repository.MysqlConfig{
		User:   os.Getenv("DBUSER"),
		Pass:   os.Getenv("DBPASS"),
		Host:   os.Getenv("DBHOST"),
		Port:   os.Getenv("DBPORT"),
		DBName: os.Getenv("DBNAME"),
	}

	repo, err = repository.NewMySQLRepository(cfg)
	handleErr(err)

	server := api.New(repo)

	fmt.Printf("Server listening on port: %s\n", port)

	err = http.ListenAndServe(":"+port, server)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
