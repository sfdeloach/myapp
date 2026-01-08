package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"server/helpers"
	"server/models"
)

type ContactHandler struct {
	DB    *gorm.DB
	Store *session.Store
}

func NewContactHandler(db *gorm.DB, store *session.Store) *ContactHandler {
	return &ContactHandler{
		DB:    db,
		Store: store,
	}
}

func (h *ContactHandler) Index(c *fiber.Ctx) error {
	// Get session and flash message if available
	sess, err := h.Store.Get(c)
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
	query := h.DB.Model(&models.Contact{})

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
}

func (h *ContactHandler) New(c *fiber.Ctx) error {
	return c.Render("form", fiber.Map{
		"FormAction": "/contacts/new",
		"View":       "new"}, "layouts/main")
}

func (h *ContactHandler) Create(c *fiber.Ctx) error {
	contact := new(models.Contact)

	if err := c.BodyParser(contact); err != nil {
		return c.Status(400).SendString("Invalid request data")
	}

	// Validate and sanitize data
	if err := helpers.ValidateContact(contact, h.DB); err != nil {
		return c.Render("form", fiber.Map{
			"Contact":    contact,
			"FormAction": "/contacts/new",
			"View":       "new",
			"Error":      err.Error()}, "layouts/main")
	}

	// Save to database
	if err := h.DB.Create(contact).Error; err != nil {
		return c.Status(500).SendString("Failed to save contact.")
	}

	// Get session
	sess, err := h.Store.Get(c)
	if err != nil {
		return err
	}

	// Set flash message
	sess.Set("flash_success", "Contact created.")
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/contacts")
}

func (h *ContactHandler) Show(c *fiber.Ctx) error {
	// Get session and flash message if available
	sess, err := h.Store.Get(c)
	if err != nil {
		return err
	}
	flashMessage := sess.Get("flash_success")

	// Delete the flash if it exists
	if flashMessage != nil {
		sess.Delete("flash_success")
		sess.Save()
	}

	var contact models.Contact

	if err := h.DB.Where("id = ?", c.Params("contactID")).First(&contact).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).Render("not-found", nil, "layouts/main")
		}
		return c.Status(500).SendString("Database error")
	}

	return c.Render("show", fiber.Map{"Contact": contact, "Flash": flashMessage}, "layouts/main")
}

func (h *ContactHandler) Edit(c *fiber.Ctx) error {
	var contact models.Contact

	if err := h.DB.Where("id = ?", c.Params("contactID")).First(&contact).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).Render("not-found", nil, "layouts/main")
		}
		return c.Status(500).SendString("Database error")
	}

	return c.Render("form", fiber.Map{
		"Contact":    contact,
		"FormAction": fmt.Sprintf("/contacts/%d/edit", contact.ID),
		"View":       "edit"}, "layouts/main")
}

func (h *ContactHandler) Update(c *fiber.Ctx) error {
	contact := new(models.Contact)

	if err := c.BodyParser(contact); err != nil {
		return c.Status(400).SendString("Invalid request data")
	}

	// Set the ID from the URL parameter
	id, err := strconv.ParseUint(c.Params("contactID"), 10, 32)
	if err != nil {
		return c.Status(400).SendString("Invalid contact ID")
	}
	contact.ID = uint(id)

	// Validate and sanitize data
	if err := helpers.ValidateContact(contact, h.DB); err != nil {
		return c.Render("form", fiber.Map{
			"Contact":    contact,
			"FormAction": fmt.Sprintf("/contacts/%d/edit", contact.ID),
			"View":       "edit",
			"Error":      err.Error(),
		}, "layouts/main")
	}

	// Update database
	if err := h.DB.Updates(contact).Error; err != nil {
		return c.Status(500).SendString("Failed to update contact.")
	}

	// Get session
	sess, err := h.Store.Get(c)
	if err != nil {
		return err
	}

	// Set flash message
	sess.Set("flash_success", "Contact updated.")
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect(fmt.Sprintf("/contacts/%d", contact.ID))
}

func (h *ContactHandler) Delete(c *fiber.Ctx) error {
	// Parse the ID from the URL parameter
	id, err := strconv.ParseUint(c.Params("contactID"), 10, 32)
	if err != nil {
		return c.Status(400).SendString("Invalid contact ID")
	}

	// Delete contact from database by ID
	if err := h.DB.Delete(&models.Contact{}, uint(id)).Error; err != nil {
		return c.Status(500).SendString("Failed to delete contact.")
	}

	// Get session
	sess, err := h.Store.Get(c)
	if err != nil {
		return err
	}

	// Set and save the flash message to the session
	sess.Set("flash_success", "Contact deleted.")
	if err := sess.Save(); err != nil {
		return err
	}

	// Set status to override DELETE request to /contacts
	return c.Redirect("/contacts", fiber.StatusSeeOther)
}

func (h *ContactHandler) ValidateEmail(c *fiber.Ctx) error {
	email := c.Query("email")

	// Parse contact ID from URL parameter (0 for new contacts)
	var contactID uint
	if idParam := c.Params("contactID"); idParam != "" {
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			return c.Status(400).SendString("Invalid contact ID")
		}
		contactID = uint(id)
	}

	if err := helpers.ValidateEmail(email, h.DB, contactID); err != nil {
		return c.SendString(err.Error())
	}

	return c.SendString("")
}
