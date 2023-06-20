package routes

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hussain4real/simple-url-shortener/database"
	"github.com/hussain4real/simple-url-shortener/helpers"
	"os"
	"strconv"
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
	r2 := database.CreateClient(1)
	defer func(r2 *redis.Client) {
		err := r2.Close()
		if err != nil {
			panic(err)
		}
	}(r2)
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		//val, _ = r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}
	//check if url is valid
	if !govalidator.IsURL(body.URL) {
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

	var id string

	if body.CustomShortURL == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShortURL
	}

	r := database.CreateClient(0)
	defer func(r *redis.Client) {
		err := r.Close()
		if err != nil {
			panic(err)
		}
	}(r)

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom URL already exists",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot connect to server",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShortURL:  "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = int(ttl / time.Nanosecond / time.Minute)

	resp.CustomShortURL = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
