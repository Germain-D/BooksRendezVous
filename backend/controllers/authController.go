// controllers/authController.go

package controllers

import (
	"booksrendezvous-backend/database"
	"booksrendezvous-backend/models"
	"booksrendezvous-backend/utils"
	"errors"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var config, _ = utils.LoadConfig()

// Initialize logger with the log level from environment variables
var sugar = utils.SugaredLogger

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func Login(c *fiber.Ctx) error {

	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type UserData struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password []byte `json:"password"`
	}

	input := new(LoginInput)
	var userData UserData
	usermodels, err := new(models.User), *new(error)

	if err := c.BodyParser(&input); err != nil {
		sugar.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	if isEmail(input.Email) {
		usermodels, err = getUserByEmail(input.Email)
		if err != nil {
			sugar.Errorw("Failed to get user by email", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
		}
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	} else if usermodels == nil {
		CheckPasswordHash([]byte(input.Password), []byte(""))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	} else {
		userData = UserData{
			ID:       usermodels.ID,
			Name:     usermodels.Name,
			Email:    usermodels.Email,
			Password: usermodels.Password,
		}
	}

	if !CheckPasswordHash([]byte(input.Password), userData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": nil})
	}

	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userData.Name
	claims["user_id"] = userData.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "token": t, "pseudo": userData.Name, "uuid": userData.ID, "email": userData.Email})

}

// controllers/authController.go

func Register(c *fiber.Ctx) error {
	sugar.Info("Received a register request")

	// Load configuration to get authorized emails
	config, err := utils.LoadConfig()
	if err != nil {
		sugar.Errorw("Failed to load config", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Parse request body
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		sugar.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Check if the email is in the authorized emails list
	if !utils.Contains(config.AuthorizedEmails, data["email"]) {
		sugar.Warnw("Email not authorized", "email", data["email"])
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email not authorized",
		})
	}

	// Check if the email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", data["email"]).First(&existingUser).Error; err == nil {
		sugar.Warnw("Email already exists", "email", data["email"])
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		sugar.Errorw("Failed to hash password", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create new user
	user := models.User{
		Name:     data["username"],
		Email:    data["email"],
		Password: hashedPassword,
	}

	// Insert user into database
	if err := database.DB.Create(&user).Error; err != nil {
		sugar.Errorw("Failed to create user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Get the user ID
	var userID string
	database.DB.Table("users").Select("id").Where("email = ?", data["email"]).Scan(&userID)

	publicusers := models.Publicusers{
		UserID:   user.ID,
		IsPublic: false,
	}

	// Insert public user into database
	if err := database.DB.Create(&publicusers).Error; err != nil {
		sugar.Errorw("Failed to create public user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create public user",
		})
	}

	// Return success response
	sugar.Infow("User registered successfully", "email", data["email"])
	return c.JSON(fiber.Map{
		"message": "User registered successfully",
	})

}

func User(c *fiber.Ctx) error {

	sugar.Info("Received a user request")

	uiidStr, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := uuid.Parse(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to parse uuid", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse uuid",
		})
	}

	var user models.User

	// Query user from database using ID
	database.DB.Where("id =?", userID).First(&user)

	// Return user details as JSON response
	sugar.Infow("User details retrieved successfully", "email", user.Email)
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	sugar.Info("Received a logout request")

	_, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Return success response indicating logout was successful
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func PasswordChange(c *fiber.Ctx) error {
	sugar.Info("Received a reset password request")

	uiidStr, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := uuid.Parse(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to parse uuid", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse uuid",
		})
	}

	var user models.User

	// Query user from database using ID
	database.DB.Where("id =?", userID).First(&user)

	// Parse request body
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		sugar.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if !CheckPasswordHash([]byte(data["oldpassword"]), user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": nil})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["newpassword"]), bcrypt.DefaultCost)
	if err != nil {
		sugar.Errorw("Failed to hash password", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Update user password in database
	database.DB.Model(&user).Update("password", hashedPassword)

	// Return success response
	sugar.Infow("Password reset successfully")
	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}
