package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eoliveira-daitan/go-challenges/internal/api"
	"github.com/eoliveira-daitan/go-challenges/internal/repository"
	"github.com/joho/godotenv"
)

const (
	ORM     = "ORM"
	VANILLA = "VANILLA"
)

func initRepository(dbImpl string, cfg repository.DBConfig) (repository.Repository, error) {
	switch dbImpl {
	case VANILLA:
		return repository.NewMySQLRepository(cfg)
	case ORM:
		return repository.NewOrmRepository(cfg)
	}

	return nil, fmt.Errorf("%q is not a valid option for env var DBIMPL.\nValid options (ORM and VANILLA).\nExample: DBIMPL=VANILLA", dbImpl)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Create a .env file in the root directory. Use `.env.example` as a start point")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbImpl := strings.ToUpper(os.Getenv("DBIMPL"))
	if dbImpl == "" {
		dbImpl = VANILLA
	}

	bearerToken := "Bearer " + os.Getenv("BEARER_TOKEN")
	auth := func(token string) bool {
		return bearerToken == token
	}

	cfg := repository.DBConfig{
		User:   os.Getenv("DBUSER"),
		Pass:   os.Getenv("DBPASS"),
		Host:   os.Getenv("DBHOST"),
		Port:   os.Getenv("DBPORT"),
		DBName: os.Getenv("DBNAME"),
	}

	var repo repository.Repository
	repo, err = initRepository(dbImpl, cfg)
	handleErr(err)

	server := api.New(repo, auth)

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
