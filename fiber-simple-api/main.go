package main

import (
	"first/book"
	"first/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func welcome(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello from Fiber 😂")
}

func initDatabase() {
	var err error
	database.DBConnection, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	log.Println("Connection Opened to database")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", welcome)

	app.Get("/api/v1/book", book.GetBooks)
	app.Get("/api/v1/book/:id", book.GetBook)
	app.Post("/api/v1/book/", book.NewBook)
	app.Delete("/api/v1/book/:id", book.DeleteBook)
}

func main() {
	app := fiber.New()

	initDatabase()
	setupRoutes(app)

	log.Fatal(app.Listen(":5252"))
}
