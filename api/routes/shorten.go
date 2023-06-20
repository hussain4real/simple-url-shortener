package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hussain4real/simple-url-shortener/helpers"
	"time"
)

type request struct {
	URL            string        `json:"url"`
	CustomShortURL string        `json:"short_url"`
	Expiry         time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShortURL  string        `json:"short_url"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_remaining"`
	XRateLimitReset int           `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	//implement rate limiting

	//check if url is valid
	if !goValidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is invalid",
		})
	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "URL is not allowed",
		})
	}
	//enforce https, SSL
	body.URL = helpers.EnforceHTTP(body.URL)
}
