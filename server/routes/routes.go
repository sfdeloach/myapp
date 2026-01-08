package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"server/handlers"
)

func Setup(app *fiber.App, db *gorm.DB, store *session.Store) {
	// Initialize handlers
	contactHandler := handlers.NewContactHandler(db, store)

	// Root redirect
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("contacts")
	})

	// Contact routes
	app.Get("/contacts", contactHandler.Index)
	app.Get("/contacts/new", contactHandler.New)
	app.Post("/contacts/new", contactHandler.Create)
	app.Get("/contacts/:contactID", contactHandler.Show)
	app.Get("/contacts/:contactID/edit", contactHandler.Edit)
	app.Post("/contacts/:contactID/edit", contactHandler.Update)
	app.Delete("/contacts/:contactID/", contactHandler.Delete)

	// Form validation routes
	app.Get("/contacts/:contactID/validate/email", contactHandler.ValidateEmail)

	// 404 handler - must be last
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("not-found", nil, "layouts/main")
	})
}
