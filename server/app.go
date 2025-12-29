package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type Contact struct {
	gorm.Model        // provides ID, CreatedAt, UpdatedAt, and DeletedAt
	First      string `gorm:"size:50"`
	Last       string `gorm:"size:50"`
	Phone      string `gorm:"size:50"`
	Email      string `gorm:"not null"`
}

func main() {
	// Get values from environment variables
	host, hostOk := os.LookupEnv("POSTGRES_HOST")
	port, portOk := os.LookupEnv("POSTGRES_PORT")
	user, userOk := os.LookupEnv("POSTGRES_USER")
	password, passwordOk := os.LookupEnv("POSTGRES_PASSWORD")
	dbname, dbnameOk := os.LookupEnv("POSTGRES_DB")
	sslmode, sslmodeOk := os.LookupEnv("POSTGRES_SSLMODE")

	// Basic validation (add more as needed)
	if !hostOk || !portOk || !userOk || !passwordOk || !dbnameOk || !sslmodeOk {
		log.Fatal("Missing required DB environment variables")
	}

	// Build DSN dynamically
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successful database connection")
	}

	// Auto-migrate the User model
	db.AutoMigrate(&Contact{})

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("contacts")
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		// dummy data
		contacts := []map[string]string{
			{"first": "Alan", "last": "Adams", "phone": "321-555-1598", "email": "aadams@xyz.com"},
			{"first": "Bob", "last": "Bryant", "phone": "321-555-1745", "email": "bbryant@xyz.com"},
			{"first": "Chuck", "last": "Connors", "phone": "321-555-1652", "email": "cconnors@xyz.com"},
			{"first": "Dave", "last": "Denver", "phone": "321-555-1598", "email": "ddenver@xyz.com"},
		}

		contacts_searched := []map[string]string{
			{"first": "Alan", "last": "Adams", "phone": "321-555-1598", "email": "aadams@xyz.com"},
		}

		search := c.Query("q")
		if search != "" {
			// implement database function here to search?
			return c.Render("index", fiber.Map{"Contacts": contacts_searched})
		} else {
			return c.Render("index", fiber.Map{"Contacts": contacts})
		}
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
