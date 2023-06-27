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
	// get auth user id
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	// get all shortlies by user id
	shortlies, err := models.GetAllShortliesByUserId(*userID)
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
	// get auth user id
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	// get shortly id from params
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id " + err.Error(),
		})
	}
	// get shortly by id
	shortly, err := models.GetShortlyById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not get shortly " + err.Error(),
		})
	}
	// check if user is owner of shortly
	if shortly.UserID != *userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    shortly,
	})
}

// create a shortly
func CreateShortly(c *fiber.Ctx) error {
	// check if there is a user id in the token
	userID, err := utils.GetUserIDFromToken(c)
	//if user id is not in the token, create a shortly url omitting the user id
	if err != nil {

		//create a guest user
		var guestUser *models.User

		guestUser, err = models.CreateGuestUser()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "could not create guest user " + err.Error(),
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

		// Overide the random value to false since it is a guest user
		shortly.Random = true
		// Generate a random URL for the shortly

		shortly.ShortURL = utils.RandomUrl(8)

		// Set the guest user ID on the shortly record
		shortly.UserID = guestUser.UserID

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

	//if user id is in the token, create a shortly url and save it to the database

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
