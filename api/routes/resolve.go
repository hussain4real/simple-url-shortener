package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/hussain4real/simple-url-shortener/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	r := database.CreateClient(0)
	defer func(r *redis.Client) {
		err := r.Close()
		if err != nil {
			panic(err)
		}
	}(r)
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot connect to database",
		})
	}
	rInr := database.CreateClient(1)
	defer func(rInr *redis.Client) {
		err := rInr.Close()
		if err != nil {
			panic(err)
		}
	}(rInr)
	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, fiber.StatusMovedPermanently)
}
