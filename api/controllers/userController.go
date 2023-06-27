package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hussain4real/simple-url-shortener/models"
	"golang.org/x/crypto/bcrypt"
)

// Home func
func Home(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

// GetUser func
func GetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not get user",
		})
	}

	var user models.User

	db := models.DB
	db.Preload("Shortlies").Where("id = ?", (*claims)["id"]).First(&user)

	return c.JSON(user)
}

// Logic to register user
func RegisterUser(c *fiber.Ctx) error {
	db := models.DB
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14) //GenerateFromPassword returns the bcrypt hash of the password at the given cost i.e. (14 in our case).

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		UserName:  data["user_name"],
		Email:     data["email"],
		Password:  string(password),
	}
	db.Create(&user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
		"user":    user,
	})

}

const SecretKey = "secret"

// LoginUser func
func LoginUser(c *fiber.Ctx) error {
	db := models.DB
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	db.Where("email = ?", data["email"]).First(&user) //Check the email is present in the DB

	if user.ID == 0 { //If the ID return is '0' then there is no such email present in the DB
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil { //Compare the password from the request with the hash from the DB
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	err2, done := Auth(c, user)
	if done {
		return err2
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"user":    user,
	})
}

func Auth(c *fiber.Ctx, user models.User) (error, bool) {
	//Create JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), //Set the expiration time of the token to 24 hours
	})
	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		}), true
	}

	//Set the cookie in the browser
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24), //Set the expiration time of the cookie to 24 hours
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return nil, false
}

// LogoutUser func
func LogoutUser(c *fiber.Ctx) error {
	//Remove the cookie from the browser
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), //Set the expiration time of the cookie to 24 hours
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

// UpdateUser func
func UpdateUser(c *fiber.Ctx) error {
	return c.SendString("Update User")
}

// DeleteUser func
func DeleteUser(c *fiber.Ctx) error {
	return c.SendString("Delete User")
}
