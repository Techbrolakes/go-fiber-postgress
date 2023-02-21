package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := new(Book)

	if err := c.BodyParser(book); err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "request failed"})
		return err
	}

	if err := r.DB.Create(book).Error; err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "Could not create book"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book has been added"})
	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	if err := r.DB.Find(bookModels).Error; err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Could not get books"})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books Fetched Successfully", "data": bookModels})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/books", r.CreateBook)
	api.Delete("/books/:id", r.DeleteBook)
	api.Get("/books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not establish a database connection")
	}

	r := &Repository{DB: db}

	app := fiber.New()
	r.SetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
