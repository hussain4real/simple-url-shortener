package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hussain4real/simple-url-shortener/models"
)

// get all redirects
func GetAllRedirects(c *fiber.Ctx) error {
	shortlies, err := models.GetAllShortlies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not get all redirects " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    shortlies,
	})
}
