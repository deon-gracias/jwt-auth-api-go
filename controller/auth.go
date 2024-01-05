package controller

import (
	"auth-api-go/config"
	"auth-api-go/database"
	"auth-api-go/model"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckPassowrdHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserByEmail(email string) (*model.User, error) {
	db := database.DB
	var user model.User

	if err := db.Where(&model.User{Email: email}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error on login request",
			"data":    err,
		})
	}

	// userModel, err := new(model.User), *new(error)

	userModel, err := GetUserByEmail(input.Email)

	if userModel == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err,
		})
	}

	if !CheckPassowrdHash(input.Password, userModel.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Password",
			"data":    nil,
		})
	}

	// claims := token.Claims.(jwt.MapClaims)
	claims := jwt.MapClaims{
		"email": userModel.Email,
		"id":    userModel.ID.String(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Expires after 72 hours
	}

	// token := jwt.New(jwt.SigningMethodHS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.Config("SECRET_KEY")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully Logged In",
		"data": fiber.Map{
			"user":  userModel,
			"token": t,
		},
	})
}
