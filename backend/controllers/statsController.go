package controllers

import (
	"booksrendezvous-backend/database"
	"booksrendezvous-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type statstoreturn struct {
	TotalBooks     int     `json:"totalBooks"`
	CompletedBooks int     `json:"completedBooks"`
	ToReadBooks    int     `json:"toReadBooks"`
	ReadingBooks   int     `json:"readingBooks"`
	FavoriteBooks  int     `json:"favoriteBooks"`
	TotalPages     int     `json:"totalPages"`
	AverageRating  float64 `json:"averageRating"`
}

func GetStats(c *fiber.Ctx) error {

	sugar.Info("Received a user stats request")

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

	var userstats models.UserStat
	// Query user stats from database using ID
	database.DB.Where("user_id =?", userID).First(&userstats)

	// Check if user stats exist
	if userstats.UserID == "" {
		// User stats do not exist, compute from scratch
		userstats = ComputeStatsFromScratch(userID)
		// Save user stats to database
		database.DB.Create(&userstats)
	} else {
		// User stats exist, return as JSON response
		sugar.Infow("User stats retrieved successfully")
		// parse userstats to statstoreturn
		toreturn := statstoreturn{
			TotalBooks:     userstats.TotalBooks,
			CompletedBooks: userstats.CompletedBooks,
			ToReadBooks:    userstats.ToReadBooks,
			ReadingBooks:   userstats.ReadingBooks,
			FavoriteBooks:  userstats.FavoriteBooks,
			TotalPages:     userstats.TotalPages,
			AverageRating:  userstats.AverageRating,
		}

		return c.JSON(toreturn)
	}

	// Return user stats as JSON response
	sugar.Infow("User stats retrieved successfully")
	toreturn := statstoreturn{
		TotalBooks:     userstats.TotalBooks,
		CompletedBooks: userstats.CompletedBooks,
		ToReadBooks:    userstats.ToReadBooks,
		ReadingBooks:   userstats.ReadingBooks,
		FavoriteBooks:  userstats.FavoriteBooks,
		TotalPages:     userstats.TotalPages,
		AverageRating:  userstats.AverageRating,
	}

	return c.JSON(toreturn)
}

func ComputeStatsFromScratch(userID uuid.UUID) models.UserStat {
	// Query all books from database
	var books []models.Book
	database.DB.Where("user_id = ?", userID).Find(&books)

	// Compute total books
	totalBooks := len(books)

	// Compute total pages/ totalfavoritebooks
	totalPages := 0
	totalFavoriteBooks := 0
	completedBooks := 0
	toreadBooks := 0
	readingBooks := 0
	totalRating := 0
	for _, book := range books {
		totalPages += book.PageCount
		totalRating += book.Rating
		if book.Favorite {
			totalFavoriteBooks++
		}
		if book.Status == "finished" {
			completedBooks++
		}
		if book.Status == "to-read" {
			toreadBooks++
		}
		if book.Status == "reading" {
			readingBooks++
		}

	}
	var userstats models.UserStat
	// Compute average rating
	if totalBooks > 0 {
		averageRating := float64(totalRating) / float64(totalBooks)

		userstats.AverageRating = averageRating
	}

	userstats.UserID = userID.String()
	userstats.TotalBooks = totalBooks
	userstats.TotalPages = totalPages
	userstats.FavoriteBooks = totalFavoriteBooks
	userstats.CompletedBooks = completedBooks
	userstats.ToReadBooks = toreadBooks
	userstats.ReadingBooks = readingBooks

	return userstats
}

func OnAddUpdateStats(userID uuid.UUID, book models.Book) {
	// Query user stats from database using ID
	var userstats models.UserStat
	database.DB.Where("user_id =?", userID).First(&userstats)

	// Update user stats
	userstats.TotalBooks++
	userstats.TotalPages += book.PageCount
	userstats.AverageRating = (float64(userstats.TotalBooks-1)*userstats.AverageRating + float64(book.Rating)) / float64(userstats.TotalBooks)
	if book.Favorite {
		userstats.FavoriteBooks++
	}
	if book.Status == "finished" {
		userstats.CompletedBooks++
	}
	if book.Status == "to-read" {
		userstats.ToReadBooks++
	}
	if book.Status == "reading" {
		userstats.ReadingBooks++
	}

	// Save updated user stats to database
	database.DB.Save(&userstats)
}

func OnDeleteUpdateStats(userID uuid.UUID, book models.Book) {
	// Query user stats from database using ID
	var userstats models.UserStat
	database.DB.Where("user_id =?", userID).First(&userstats)

	// Update user stats
	userstats.TotalBooks--
	userstats.TotalPages -= book.PageCount
	userstats.AverageRating = (float64(userstats.TotalBooks+1)*userstats.AverageRating - float64(book.Rating)) / float64(userstats.TotalBooks)
	if book.Favorite {
		userstats.FavoriteBooks--
	}
	if book.Status == "finished" {
		userstats.CompletedBooks--
	}
	if book.Status == "to-read" {
		userstats.ToReadBooks--
	}
	if book.Status == "reading" {
		userstats.ReadingBooks--
	}

	// Save updated user stats to database
	database.DB.Save(&userstats)
}

func OnChangeUpdateStats(userID uuid.UUID, newbook models.Book, oldbook models.Book) {
	// Query user stats from database using ID
	var userstats models.UserStat
	database.DB.Where("user_id =?", userID).First(&userstats)

	// Action pour le changement de statut
	if newbook.Status != oldbook.Status {
		if newbook.Status == "finished" {
			userstats.CompletedBooks++
		} else if newbook.Status == "to-read" {
			userstats.ToReadBooks++
		} else if newbook.Status == "reading" {
			userstats.ReadingBooks++
		}

		if oldbook.Status == "finished" {
			userstats.CompletedBooks--
		} else if oldbook.Status == "to-read" {
			userstats.ToReadBooks--
		} else if oldbook.Status == "reading" {
			userstats.ReadingBooks--
		}
	}
	if newbook.Favorite != oldbook.Favorite {
		if newbook.Favorite {
			userstats.FavoriteBooks++
		} else {
			userstats.FavoriteBooks--
		}
	}

	if newbook.Rating != oldbook.Rating {
		userstats.AverageRating = (float64(userstats.TotalBooks)*userstats.AverageRating - float64(oldbook.Rating) + float64(newbook.Rating)) / float64(userstats.TotalBooks)
	}

	// Sauvegarder les stats mises à jour
	database.DB.Save(&userstats)
}
