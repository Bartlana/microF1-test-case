package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"MicroF1-test-case/user"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	sqlConnection := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	fmt.Println("sqlConnection", sqlConnection)
	db, err := sql.Open("postgres", sqlConnection)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("failed to ping the database: %v\n", err)
	}

	log.Println("Successfully connected to the database!")

	userService := user.NewService(db)
	h := user.NewHandler(userService)

	app := fiber.New()

	app.Get("/users", h.GetUsers)
	app.Get("/users/:id", h.GetUser)
	app.Post("/users", h.CreateUser)
	app.Put("/users/:id", h.UpdateUser)
	app.Delete("/users/:id", h.DeleteUser)

	log.Fatal(app.Listen(":8080"))
}
