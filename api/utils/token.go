package utils

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const SecretKey = "secret"

func GetUserIDFromToken(c *fiber.Ctx) (*uint, error) {
	// Get the JWT token from the request header
	cookie := c.Cookies("jwt")

	// Parse the JWT token and extract the user ID
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, err
	}

	userID, err := strconv.Atoi(fmt.Sprintf("%.0f", (*claims)["id"]))
	if err != nil {
		return nil, err
	}

	userIDUint := uint(userID)
	return &userIDUint, nil
}
