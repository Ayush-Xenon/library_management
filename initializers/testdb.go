package initializers

import (
	"database/sql"
	"fmt"
	"library_management/models"

	_ "github.com/mattn/go-sqlite3" // Importing SQLite driver for in-memory database
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	// Connect to the test database
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("not connetceknfed")
		
		//panic("failed to connect to the test database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Library{}, &models.Book{}, &models.RequestEvent{}, &models.IssueRegistry{})

	return db
}
func MockDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	// Initialize schema and seed data if necessary
	return db
}
func CloseTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance")
	}
	err = sqlDB.Close()
	if err != nil {
		panic("failed to close the database")
	}
}
