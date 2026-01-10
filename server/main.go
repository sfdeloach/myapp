package main

import (
	"flag"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"

	"server/database"
	"server/helpers"
	"server/routes"
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

	// Html template engine for views
	engine := html.New("./views", ".html")

	// Register custom template function(s)
	engine.AddFunc("formatPhone", helpers.FormatPhone)
	engine.AddFunc("getPage", helpers.GetPage)

	// Set template engine
	app := fiber.New(fiber.Config{Views: engine})

	// Serve static files
	app.Static("/", "./public")

	// Create a session store for flash messages
	store := session.New()

	// Setup routes
	routes.Setup(app, db, store)

	log.Fatal(app.Listen(":3000"))
}
