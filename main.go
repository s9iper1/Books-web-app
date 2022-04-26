package main

import (
	"books-app/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBooks(c *fiber.Ctx) error {
	book := Book{}
	error := c.BodyParser(&book)

	if error != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return error
	}
	err := r.DB.Create(&book).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not create the book"})
		return err
	}
	c.Status(http.StatusCreated).JSON(&fiber.Map{"message": "Book has been created"})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get the books"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "books fetched successfully",
		"data": bookModels})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "id cannot be found"})
		return nil
	}
	error := r.DB.Delete(bookModel, id)
	if error.Error != nil {
		context.Status(http.StatusConflict).JSON(&fiber.Map{"message": "could not delete the book"})
		return error.Error
	}
	context.Status(http.StatusNoContent).JSON(&fiber.Map{"message": "successfully deleted"})
	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}
	fmt.Println("id is ", id)
	error := r.DB.Where("id = ?", id).First(bookModel).Error
	if error != nil {
		context.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "could not get the book"})
		return error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Success",
		"data": bookModel})
	return nil
}

func (r *Repository) SetUpRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBooks)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookById)
	api.Get("/books", r.GetBooks)
}

func main() {
	error := godotenv.Load(".env")

	if error != nil {
		log.Fatal(error)

	}

	//config := &storage.Config{
	//	Host:     os.Getenv("DB_HOST"),
	//	Port:     os.Getenv("DB_PORT"),
	//	Password: os.Getenv("DB_PASSWORD"),
	//	User:     os.Getenv("DB_USER"),
	//	SSLMode:  os.Getenv("DB_SSLMODE"),
	//	DBName:   os.Getenv("DB_NAME"),
	//}
	//
	//db, error := storage.NewConnection(config)
	//
	//if error != nil {
	//	log.Fatal("could not load the database")
	//}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate the database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetUpRoutes(app)
	app.Listen(":3000")
}
