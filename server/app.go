package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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

	// Register custom template function
	engine := html.New("./views", ".html")
	engine.AddFunc("formatPhone", helpers.FormatPhone)
	app := fiber.New(fiber.Config{Views: engine})

	// Serve static files
	app.Static("/", "./public")

	// Create a session store for flash messages
	store := session.New()

	/*
	 *  Routes
	 */

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("contacts")
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		// Get session and flash message if available
		sess, err := store.Get(c)
		if err != nil {
			return err
		}
		flashMessage := sess.Get("flash_success")

		// Delete the flash if it exists
		if flashMessage != nil {
			sess.Delete("flash_success")
			sess.Save()
		}

		var contacts []models.Contact
		query := db.Model(&models.Contact{})

		// Check for query parameters in the URL
		if q := c.Query("q"); q != "" {
			search := "%" + q + "%"
			query = query.Where("first ILIKE ? OR last ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR CONCAT(first, ' ', last) ILIKE ?",
				search, search, search, search, search)
		}

		if err := query.Find(&contacts).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Render("index",
			fiber.Map{"Contacts": contacts, "QueryParam": c.Query("q"), "Flash": flashMessage},
			"layouts/main")
	})

	app.Get("/contacts/new", func(c *fiber.Ctx) error {
		return c.Render("new", nil, "layouts/main")
	})

	app.Post("/contacts/new", func(c *fiber.Ctx) error {
		contact := new(models.Contact)

		if err := c.BodyParser(contact); err != nil {
			return c.Status(400).SendString("Invalid request data")
		}

		// Validate and sanitize data
		if err := helpers.ValidateContact(contact); err != nil {
			return c.Status(422).Render("new", fiber.Map{
				"Contact": contact,
				"Error":   err.Error(),
			}, "layouts/main")
		}

		// Save to database
		if err := db.Create(contact).Error; err != nil {
			// Check for unique constraint violation on email
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
				return c.Status(422).Render("new", fiber.Map{
					"Contact": contact,
					"Error":   "a contact with this email already exists",
				}, "layouts/main")
			}
			return c.Status(500).SendString("Failed to save contact. Please try again.")
		}

		// Get session
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		// Set flash message
		sess.Set("flash_success", "Contact created successfully!")
		if err := sess.Save(); err != nil {
			return err
		}

		return c.Redirect("/contacts")
	})

	log.Fatal(app.Listen(":3000"))
}
