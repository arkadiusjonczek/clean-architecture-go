package mysql

import (
	"testing"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
)

// TestMySQLBasketRepository demonstrates how to use the MySQL basket repository
// This test requires a running MySQL database
func TestMySQLBasketRepository(t *testing.T) {
	// Skip this test if no database is available
	t.Skip("Skipping integration test - requires MySQL database")

	// Configure database connection
	config := DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		Username: "root",
		Password: "password",
		Database: "ecommerce",
	}

	// Create database connection
	db, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
	}()

	// Create repository
	repository := NewMySQLBasketRepository(db)

	// Create a new basket
	basketFactory := entities.NewBasketFactory()
	basket, err := basketFactory.NewBasket("user123")
	if err != nil {
		t.Fatalf("Failed to create basket: %v", err)
	}

	// Add some items
	basket.AddItem("product1", 2)
	basket.AddItem("product2", 1)

	// Save basket
	basketID, err := repository.Save(basket)
	if err != nil {
		t.Fatalf("Failed to save basket: %v", err)
	}

	// Find basket by ID
	foundBasket, err := repository.Find(basketID)
	if err != nil {
		t.Fatalf("Failed to find basket: %v", err)
	}

	// Verify basket data
	if foundBasket.GetUserID() != "user123" {
		t.Errorf("Expected user ID 'user123', got '%s'", foundBasket.GetUserID())
	}

	// Verify items
	items := foundBasket.GetItems()
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Find basket by user ID
	userBasket, err := repository.FindByUserId("user123")
	if err != nil {
		t.Fatalf("Failed to find basket by user ID: %v", err)
	}

	if userBasket.GetID() != basketID {
		t.Errorf("Expected basket ID '%s', got '%s'", basketID, userBasket.GetID())
	}
}
