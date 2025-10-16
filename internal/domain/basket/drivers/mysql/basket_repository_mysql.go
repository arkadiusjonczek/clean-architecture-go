package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
)

var _ entities.BasketRepository = (*MySQLBasketRepository)(nil)

type MySQLBasketRepository struct {
	db *sql.DB
}

func NewMySQLBasketRepository(db *sql.DB) entities.BasketRepository {
	return &MySQLBasketRepository{
		db: db,
	}
}

// BasketDB represents the basket structure in the database
type BasketDB struct {
	ID     string `json:"id"`
	UserID string `json:"userid"`
	Items  string `json:"items"` // JSON string of basket items
}

// BasketItemDB represents a basket item structure for JSON serialization
type BasketItemDB struct {
	ProductID string `json:"product_id"`
	Count     int    `json:"count"`
}

func (repository *MySQLBasketRepository) Save(basket *entities.Basket) (string, error) {
	if basket == nil {
		return "", fmt.Errorf("basket is nil")
	}

	if basket.GetID() == "" {
		basket.SetID(uuid.NewString())
	}

	// Convert basket items to JSON
	itemsJSON, err := repository.serializeItems(basket.GetItems())
	if err != nil {
		return "", fmt.Errorf("failed to serialize basket items: %w", err)
	}

	// Use INSERT ... ON DUPLICATE KEY UPDATE for upsert behavior
	query := `
		INSERT INTO baskets (id, userid, items) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		userid = VALUES(userid), 
		items = VALUES(items)
	`

	_, err = repository.db.Exec(query, basket.GetID(), basket.GetUserID(), itemsJSON)
	if err != nil {
		return "", fmt.Errorf("failed to save basket: %w", err)
	}

	return basket.GetID(), nil
}

func (repository *MySQLBasketRepository) Find(id string) (*entities.Basket, error) {
	query := `SELECT id, userid, items FROM baskets WHERE id = ?`

	var basketDB BasketDB
	err := repository.db.QueryRow(query, id).Scan(&basketDB.ID, &basketDB.UserID, &basketDB.Items)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &entities.BasketNotFoundError{}
		}
		return nil, fmt.Errorf("failed to find basket: %w", err)
	}

	// Deserialize basket items
	items, err := repository.deserializeItems(basketDB.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize basket items: %w", err)
	}

	// Create basket entity
	basket := &entities.Basket{}
	basket.SetID(basketDB.ID)
	// Note: We need to set userID and items through reflection or add setters to the entity
	// For now, we'll create a new basket using the factory pattern
	basketFactory := entities.NewBasketFactory()
	basket, err = basketFactory.NewBasketWithID(basketDB.UserID, basketDB.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create basket entity: %w", err)
	}

	// Add items to the basket
	for productID, item := range items {
		basket.AddItem(productID, item.GetCount())
	}

	return basket, nil
}

func (repository *MySQLBasketRepository) FindByUserId(userId string) (*entities.Basket, error) {
	query := `SELECT id, userid, items FROM baskets WHERE userid = ?`

	var basketDB BasketDB
	err := repository.db.QueryRow(query, userId).Scan(&basketDB.ID, &basketDB.UserID, &basketDB.Items)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &entities.BasketNotFoundError{}
		}
		return nil, fmt.Errorf("failed to find basket by user ID: %w", err)
	}

	// Deserialize basket items
	items, err := repository.deserializeItems(basketDB.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize basket items: %w", err)
	}

	// Create basket entity
	basketFactory := entities.NewBasketFactory()
	basket, err := basketFactory.NewBasketWithID(basketDB.UserID, basketDB.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create basket entity: %w", err)
	}

	// Add items to the basket
	for productID, item := range items {
		basket.AddItem(productID, item.GetCount())
	}

	return basket, nil
}

// serializeItems converts basket items map to JSON string
func (repository *MySQLBasketRepository) serializeItems(items map[string]*entities.BasketItem) (string, error) {
	itemsDB := make(map[string]BasketItemDB)
	for productID, item := range items {
		itemsDB[productID] = BasketItemDB{
			ProductID: item.GetProductID(),
			Count:     item.GetCount(),
		}
	}

	jsonData, err := json.Marshal(itemsDB)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// deserializeItems converts JSON string back to basket items map
func (repository *MySQLBasketRepository) deserializeItems(itemsJSON string) (map[string]*entities.BasketItem, error) {
	var itemsDB map[string]BasketItemDB
	err := json.Unmarshal([]byte(itemsJSON), &itemsDB)
	if err != nil {
		return nil, err
	}

	// Since BasketItem fields are private, we need to create them through the basket
	// Create a temporary basket using the factory to ensure proper initialization
	basketFactory := entities.NewBasketFactory()
	tempBasket, err := basketFactory.NewBasket("temp-user")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary basket: %w", err)
	}

	for productID, itemDB := range itemsDB {
		tempBasket.AddItem(productID, itemDB.Count)
	}

	return tempBasket.GetItems(), nil
}
