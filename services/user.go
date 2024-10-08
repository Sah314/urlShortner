package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"time"
	"urlShortner/helpers"
	"urlShortner/models"
	"urlShortner/repository"
)

func Signup(c *fiber.Ctx) error {
	body := new(models.SignupRequest)
	err := helpers.ParseRequestBody(c, body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if body.FirstName == "" || body.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}
	if body.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	if body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}
	password, err := helpers.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot hash password",
		})
	}
	log.Debug(password)
	//TODO: Check if user already exists
	user := repository.GetUser(body.Email)
	if user != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User already exists",
		})
	}
	//TODO: Save user to database
	createdUser, err := repository.CreateUser(&models.User{Email: body.Email, FirstName: body.LastName, LastName: body.FirstName, Password: password})
	if err != nil {
		return err
	}
	token, err := helpers.GenerateJWT(body.Email, createdUser.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot generate token",
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 7), // Token expiration (1 day example)
		HTTPOnly: true,                               // Prevent access via JavaScript
		Secure:   false,                              // Set to true in production (HTTPS)
		SameSite: "Lax",                              // Helps with CSRF protection
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created",
		"user":    createdUser,
	})
}

func Login(c *fiber.Ctx) error {
	var body = new(models.LoginRequest)
	err := helpers.ParseRequestBody(c, body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if body.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	if body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}
	currentUser := repository.GetUser(body.Email)
	if currentUser == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	password := currentUser.Password
	//TODO: Get password for current user from database
	result := helpers.ComparePasswords(body.Password, password)
	if !result {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}
	token, err := helpers.GenerateJWT(body.Email, currentUser.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot generate token",
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 7), // Token expiration (1 day example)
		HTTPOnly: true,                               // Prevent access via JavaScript
		Secure:   false,                              // Set to true in production (HTTPS)
		SameSite: "Lax",                              // Helps with CSRF protection
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged in",
		"user":    currentUser,
	})
}

func GoogleSignup(c *fiber.Ctx) error {

}

func UpdateUser(c *fiber.Ctx) error {
	body := new(models.UpdateUserRequest)
	err := helpers.ParseRequestBody(c, body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	var updatedUser *models.User
	if body.FirstName != "" {
		updatedUser.FirstName = body.FirstName
	}
	if body.LastName != "" {
		updatedUser.LastName = body.LastName
	}
	if body.Email != "" {
		updatedUser.Email = body.Email
	}
	if body.Password != "" {
		password, err := helpers.HashPassword(body.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot hash password",
			})
		}
		updatedUser.Password = password
	}
	_, err = repository.UpdateUser(updatedUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated",
		"user":    updatedUser,
	})
}
