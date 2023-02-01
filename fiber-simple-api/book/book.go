package book

import (
	"first/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetBooks(ctx *fiber.Ctx) error {
	db := database.DBConnection
	var books []Book
	db.Find(&books)
	return ctx.Status(200).JSON(books)
}

func GetBook(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	db := database.DBConnection
	var book Book
	db.Find(&book, id)
	return ctx.Status(200).JSON(book)
}

func NewBook(ctx *fiber.Ctx) error {
	db := database.DBConnection

	book := new(Book)
	if err := ctx.BodyParser(book); err != nil {
		return ctx.Status(503).SendString(err.Error())
	}

	db.Create(&book)
	return ctx.Status(201).JSON(book)
}

func DeleteBook(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	db := database.DBConnection

	var book Book
	db.First(&book, id)
	if book.Title == "" {
		return ctx.Status(500).SendString("No Book Found with ID")
	}
	db.Delete(&book)
	return ctx.SendString("Book succesfully deleted")
}
