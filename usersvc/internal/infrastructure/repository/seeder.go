package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

type seeder struct {
	DB *gorm.DB
}

// InitSeeder initializes the seeder
func InitSeeder(db *gorm.DB) *seeder {
	return &seeder{DB: db}
}

// Seed inserts the users data into the database
func (s *seeder) Seed() error {
	// Define the path to the JSON file in the data directory
	filePath := filepath.Join("pkg", "data", "users.json")

	// Read the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open JSON file: %w", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read JSON file: %w", err)
	}

	var users []User
	if err := json.Unmarshal(byteValue, &users); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	// Insert the data into the database
	if err := s.DB.Save(&users).Error; err != nil {
		return fmt.Errorf("could not update or create users: %w", err)
	}
	return nil
}
