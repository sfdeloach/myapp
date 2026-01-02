package main

import (
	"flag"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"server/database"
	"server/helpers"
	"server/models"
)

func main() {
	// Define CLI flags
	seedFlag := flag.Bool("seed", false, "Seed the database with test data")
	flag.Parse()

	// Initialize database
	db := database.Init()

	// Run seeder if flag is set
	if *seedFlag {
		if err := database.Seed(db); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		log.Println("Database seeded successfully!")
		os.Exit(0)
	}

	engine := html.New("./views", ".html")

	// Register custom template function
	engine.AddFunc("formatPhone", helpers.FormatPhone)

	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("contacts")
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		var contacts []models.Contact
		query := db.Model(&models.Contact{})

		if q := c.Query("q"); q != "" {
			search := "%" + q + "%"
			query = query.Where("first ILIKE ? OR last ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR CONCAT(first, ' ', last) ILIKE ?",
				search, search, search, search, search)
		}

		if err := query.Find(&contacts).Error; err != nil {
			return c.Status(500).SendString("Database error")
		}

		return c.Render("index",
			fiber.Map{"Contacts": contacts, "QueryParam": c.Query("q")},
			"layouts/main")
	})

	app.Get("/contacts/new", func(c *fiber.Ctx) error {
		return c.Render("new", nil, "layouts/main")
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
