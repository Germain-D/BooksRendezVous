package controllers

import (
	"booksrendezvous-backend/database"
	"booksrendezvous-backend/models"
	"booksrendezvous-backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetBooks returns all books for the authenticated user
func GetBooks(c *fiber.Ctx) error {
	sugar.Info("Received a Book request")

	uiidStr, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// convert string to uuid type
	userID, err := uuid.Parse(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to parse uuid", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse uuid",
		})
	}

	// get books from book table
	var books []models.Book
	if err := database.DB.Where("user_id = ?", userID).Find(&books).Error; err != nil {
		sugar.Errorw("Failed to get books", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get books",
		})
	}

	return c.JSON(fiber.Map{
		"books": books,
	})

}

func CheckAuth(c *fiber.Ctx) (string, bool) {
	// Get user claims from context (already verified by middleware)
	use := c.Locals("user").(*jwt.Token)
	claims := use.Claims.(jwt.MapClaims)
	id := claims["user_id"].(string)

	if id == "" {
		sugar.Error("No user_id in token claims")
		return "", false
	}

	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		sugar.Errorw("Failed to parse UUID from token", "error", err)
		return "", false
	}

	// Verify user exists in database
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		sugar.Errorw("User not found in database", "userID", userID)
		return "", false
	}

	sugar.Info("User authenticated successfully", "userID", userID)
	return userID.String(), true
}

type Book struct {
	ID            string   `json:"id"`
	GoogleBooksID string   `json:"googleBooksId"`
	Status        string   `json:"status"`
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Description   string   `json:"description"`
	ImageURL      string   `json:"imageUrl"`
	Rating        float64  `json:"rating"`
	Comment       string   `json:"comment"`
	Favorite      bool     `json:"favorite"`
	PageCount     int      `json:"pageCount"`
	Genres        []string `json:"genres"`
	PublishedDate string   `json:"publishedDate"`
}

type AddBookRequest struct {
	Book Book `json:"book"`
}

func AddBook(c *fiber.Ctx) error {
	sugar.Info("Received an Add Book request")

	uiidStr, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// convert string to uuid type
	userID, err := uuid.Parse(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to parse uuid", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse uuid",
		})
	}

	var request AddBookRequest
	if err := c.BodyParser(&request); err != nil {
		sugar.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Vérifiez si l'utilisateur existe
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		sugar.Errorw("User not found", "userID", userID, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if the Book contains all necessary fields
	book := request.Book
	if book.Title == "" || book.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Book data is incomplete",
		})
	}

	// Validate the rating field
	if book.Rating < 0 || book.Rating > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rating must be between 0 and 5",
		})
	}

	// Map the request to the actual `models.Book`
	realbook := models.Book{
		ID:            book.ID,
		GoogleBooksID: book.GoogleBooksID,
		Status:        book.Status,
		Title:         book.Title,
		Authors:       book.Authors,
		Description:   book.Description,
		ImageUrl:      book.ImageURL,
		Rating:        int(book.Rating), // Cast the float64 rating into int
		Comment:       book.Comment,
		UserID:        uiidStr, // Link the book to the user
		Favorite:      book.Favorite,
		PageCount:     book.PageCount,
		Genres:        book.Genres,
		PublishedDate: book.PublishedDate,
	}

	// Save the book to the database
	result := database.DB.Create(&realbook)
	if result.Error != nil {
		sugar.Errorw("Failed to save book to database", "error", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save book to database",
		})
	}

	//update user stats
	OnAddUpdateStats(userID, realbook)

	// Check if the user has unlocked any achievements
	service := services.NewAchievementService(database.DB)

	err = service.CheckAchievements(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to check achievements", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check achievements",
		})
	}

	// Return a success response
	return c.JSON(fiber.Map{
		"message": "Book added successfully",
		"book":    realbook, // Return the saved book or parts of it
	})
}

func DeleteBook(c *fiber.Ctx) error {
	sugar.Info("Received a Delete Book request")

	uiidStr, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// convert string to uuid type
	userID, err := uuid.Parse(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to parse uuid", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse uuid",
		})
	}

	// Get book ID from params
	bookID := c.Params("id")
	if bookID == "" {
		sugar.Error("Book ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Book ID is required",
		})
	}

	// Find the book in the database
	var book models.Book
	if err := database.DB.First(&book, "id = ?", bookID).Error; err != nil {
		sugar.Errorw("Book not found", "bookID", bookID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	// Verify that the book belongs to the user
	if book.UserID != userID.String() {
		sugar.Errorw("Unauthorized book deletion attempt",
			"userID", userID,
			"bookUserID", book.UserID,
		)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to delete this book",
		})
	}

	// Delete the book from the database
	if err := database.DB.Delete(&book).Error; err != nil {
		sugar.Errorw("Failed to delete book from database", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book from database",
		})
	}

	//update user stats
	OnDeleteUpdateStats(userID, book)

	// Check if the user has unlocked any achievements
	service := services.NewAchievementService(database.DB)

	err = service.CheckAchievements(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to check achievements", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check achievements",
		})
	}

	sugar.Infow("Book deleted successfully", "bookID", bookID)

	return c.JSON(fiber.Map{
		"message": "Book deleted successfully",
	})
}

func UpdateBook(c *fiber.Ctx) error {
	sugar.Info("Received an Update Book request")

	uiidStr, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// convert string to uuid type
	userID, err := uuid.Parse(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to parse uuid", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse uuid",
		})
	}

	// Get book ID from params
	bookID := c.Params("id")
	if bookID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Book ID is required",
		})
	}

	var request AddBookRequest
	if err := c.BodyParser(&request); err != nil {
		sugar.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Vérifiez si l'utilisateur existe
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		sugar.Errorw("User not found", "userID", userID, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if the Book contains all necessary fields
	livre := request.Book
	if livre.Title == "" || livre.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Book data is incomplete",
		})
	}

	// Validate the rating field
	if livre.Rating < 0 || livre.Rating > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rating must be between 0 and 5",
		})
	}

	// Map the request to the actual `models.Book`
	reallivre := models.Book{
		Status:   livre.Status,
		Rating:   int(livre.Rating), // Cast the float64 rating into int
		Comment:  livre.Comment,
		Favorite: livre.Favorite,
	}

	// Find the existing book in the database
	var book models.Book
	if err := database.DB.First(&book, "id = ?", bookID).Error; err != nil {
		sugar.Errorw("Book not found", "bookID", bookID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	// Verify that the book belongs to the user
	if book.UserID != userID.String() {
		sugar.Errorw("Unauthorized book update attempt",
			"userID", userID,
			"bookUserID", book.UserID,
		)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to delete this book",
		})
	}

	OnChangeUpdateStats(userID, reallivre, book)

	// Check if the user has unlocked any achievements
	service := services.NewAchievementService(database.DB)

	err = service.CheckAchievements(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to check achievements", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check achievements",
		})
	}

	// Update the book with the new data
	book.Comment = reallivre.Comment
	book.Rating = reallivre.Rating
	book.Status = reallivre.Status
	book.Favorite = reallivre.Favorite

	// Save the updated book to the database
	if err := database.DB.Save(&book).Error; err != nil {
		sugar.Errorw("Failed to update book in database", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book in database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Book updated successfully",
		"book":    book,
	})
}
