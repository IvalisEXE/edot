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
	if err := s.Warehouse(); err != nil {
		return err
	}
	if err := s.WarehouseStockItem(); err != nil {
		return err
	}
	return nil
}

func (s *seeder) Warehouse() error {
	// Define the path to the JSON file in the data directory
	filePath := filepath.Join("/root", "pkg", "data", "warehouse.json")

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

	var warehouse []Warehouse
	if err := json.Unmarshal(byteValue, &warehouse); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	// Insert the data into the database
	if err := s.DB.Save(&warehouse).Error; err != nil {
		return fmt.Errorf("could not update or create warehouse: %w", err)
	}

	return nil
}

func (s *seeder) WarehouseStockItem() error {
	// Define the path to the JSON file in the data directory
	filePath := filepath.Join("/root", "pkg", "data", "stock_item.json")

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

	var warehouseStockItem []WarehouseStockItem
	if err := json.Unmarshal(byteValue, &warehouseStockItem); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	// Insert the data into the database
	if err := s.DB.Save(&warehouseStockItem).Error; err != nil {
		return fmt.Errorf("could not update or create warehouse stock item: %w", err)
	}

	return nil
}
