package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model // provides ID, CreatedAt, UpdatedAt, and DeletedAt
	First      string
	Last       string
	Phone      string
	Email      string
}

func dbInit() *gorm.DB {
	// Get values from environment variables
	host, hostOk := os.LookupEnv("POSTGRES_HOST")
	port, portOk := os.LookupEnv("POSTGRES_PORT")
	user, userOk := os.LookupEnv("POSTGRES_USER")
	password, passwordOk := os.LookupEnv("POSTGRES_PASSWORD")
	dbName, dbNameOk := os.LookupEnv("POSTGRES_DB")
	sslMode, sslModeOk := os.LookupEnv("POSTGRES_SSLMODE")

	// Basic validation (add more as needed)
	if !hostOk || !portOk || !userOk || !passwordOk || !dbNameOk || !sslModeOk {
		log.Fatal("Missing required DB environment variables")
	}

	// Build DSN dynamically
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successful database connection")
	}

	// Auto-migrate the User model
	db.AutoMigrate(&Contact{})

	return db
}

func main() {
	db := dbInit()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("contacts")
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		var contacts []Contact
		var search string
		queryParam := c.Query("q")

		if queryParam != "" {
			search = "%" + queryParam + "%"
			if err := db.Where(
				db.Where("first ILIKE ?", search).
					Or("last ILIKE ?", search).
					Or("email ILIKE ?", search).
					Or("phone ILIKE ?", search),
			).Find(&contacts).Error; err != nil {
				return c.Status(500).SendString("Database error")
			}
		} else {
			if err := db.Find(&contacts).Error; err != nil {
				return c.Status(500).SendString("Database error")
			}
		}

		return c.Render("index",
			fiber.Map{"Contacts": contacts, "QueryParam": queryParam},
			"layouts/main")
	})

	// EXAMPLE - Route to create a new user
	// app.Post("/users", func(c *fiber.Ctx) error {
	// 	user := new(User)
	// 	if err := c.BodyParser(user); err != nil {
	// 		return c.Status(400).SendString(err.Error())
	// 	}
	// 	db.Create(user)
	// 	return c.JSON(user)
	// })

	log.Fatal(app.Listen(":3000"))
}
