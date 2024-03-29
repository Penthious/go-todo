package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo/domain"
	"todo/handlers"
	"todo/postgres"

	"github.com/go-pg/pg/v9"
)

func main() {

	DB := postgres.New(&pg.Options{
		User:     "tleffew",
		Password: "postgres",
		Database: "todo_dev",
	})

	defer DB.Close()

	domainDB := domain.DB{
		UserRepo: postgres.NewUserRepo(DB),
		TodoRepo: postgres.NewTodoRepo(DB),
	}
	d := &domain.Domain{DB: domainDB}

	r := handlers.SetupRouter(d)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)

	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}

}
