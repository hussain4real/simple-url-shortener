package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hussain4real/simple-url-shortener/models"
	"github.com/hussain4real/simple-url-shortener/utils"
)

// redirect to shortly url
func Redirect(c *fiber.Ctx) error {
	shortlyUrl := c.Params("redirect")
	shortly, err := models.FindByShortlyUrl(shortlyUrl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not get shortly " + err.Error(),
		})
	}

	// grab stats
	shortly.Visits += 1
	err = models.UpdateShortly(shortly)
	if err != nil {
		fmt.Printf("could not update shortly %v", err)
	}

	return c.Redirect(shortly.RedirectURL, fiber.StatusTemporaryRedirect)

}

// get all shortlies
func GetAllShortlies(c *fiber.Ctx) error {
	shortlies, err := models.GetAllShortlies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not get all shortlies " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    shortlies,
	})
}

// get a single shortly
func GetShortly(c *fiber.Ctx) error {
	//get from params
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id " + err.Error(),
		})
	}
	shortly, err := models.GetShortlyById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not get shortly " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    shortly,
	})
}

// create a shortly
func CreateShortly(c *fiber.Ctx) error {
	// Get the current user ID from the JWT token
	userID, err := utils.GetUserIDFromToken(c)

	//check if user is unauthenticated
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})

	}

	//get from body
	var shortly models.Shortly
	err = c.BodyParser(&shortly)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "could not parse request " + err.Error(),
		})
	}
	//check if the url is valid
	if !models.IsValidURL(shortly.RedirectURL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid url",
		})
	}

	// Generate a random URL if the user did not provide a URL
	if shortly.Random {
		shortly.ShortURL = utils.RandomUrl(8)
	}

	// Log the user ID and shortly data
	fmt.Printf("userID: %v\n", userID)
	fmt.Printf("shortly: %+v\n", shortly)

	// Set the user ID of the shortly to the current user's ID
	shortly.UserID = uint(*userID)

	// Create the shortly
	err = models.CreateShortly(shortly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create shortly " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    shortly,
	})

}
