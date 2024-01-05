package controller

import (
	"auth-api-go/database"
	"auth-api-go/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	claims := t.Claims.(jwt.MapClaims)
	return claims["id"] == id
}

func validUser(id uuid.UUID, p string) bool {
	var user model.User

	database.DB.First(&user, id)

	return user.Username != ""
}

// Returns existing user
func GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"error":  fmt.Sprintf("No user with ID: %s", id),
		})
	}

	var user model.User
	database.DB.Find(&user, id)

	if user.Email == "" {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"error":  fmt.Sprintf("No user with ID: %s", id),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User found",
		"data":    user,
	})
}

// Creates new user and returns user
func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(model.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "User structure is improper",
			"data":    err,
		})
	}

	hash, err := HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't hash password",
			"data":    err,
		})
	}

	user.Password = hash

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":   "error",
			"messsage": "",
			"data":     err,
		})
	}

	if err := c.BodyParser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create user",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User created",
		"data":    user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	type RequestData struct {
		Username string `json:"username"`
	}

	var request RequestData
	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "User structure isn't correct",
			"data":    err,
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"error":  fmt.Sprintf("No user with ID: %s", id),
		})
	}

	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id.String()) {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JWT Token",
			"data":    nil,
		})
	}

	db := database.DB
	var user model.User

	db.First(&user, id)
	user.Username = request.Username

	db.Save(&user)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User updated",
		"data":    user,
	})
}

// Delete user
func DeleteUser(c *fiber.Ctx) error {
	type RequestData struct {
		Password string `json:"password"`
	}

	var request RequestData
	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Input structure is improper",
			"data":    err,
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"error":  fmt.Sprintf("No user with ID: %s", id),
		})
	}

	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id.String()) {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JWT Token",
			"data":    nil,
		})
	}

	if !validUser(id, request.Password) {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't find user",
			"data":    nil,
		})
	}

	var user model.User

	database.DB.First(&user, id)
	database.DB.Unscoped().Delete(&user, id)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User Deleted",
		"data":    user,
	})
}
