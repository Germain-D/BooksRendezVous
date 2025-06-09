// database/connection.go

package database

// Import GORM and PostgreSQL driver
import (
	"booksrendezvous-backend/models"
	"booksrendezvous-backend/utils"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect to PostgreSQL database
func ConnectDB(config *utils.Config, sugar *zap.SugaredLogger) (*gorm.DB, error) {

	sugar.Debug("Connecting to database")

	// PostgreSQL connection string format
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

	sugar.Infof("Connecting to database: %s", dsn)

	// Connect to PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Book{})
	db.AutoMigrate(&models.Publicusers{})
	db.AutoMigrate(&models.UserStat{})
	db.AutoMigrate(&models.Achievement{})
	db.AutoMigrate(&models.UserAchievement{})
	db.AutoMigrate(&models.PasswordResetToken{})

	return db, nil
}
