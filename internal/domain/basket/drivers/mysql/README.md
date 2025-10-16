# MySQL Basket Repository Adapter

This package provides a MySQL implementation of the `BasketRepository` interface for the clean architecture Go project.

## Setup

### 1. Database Setup

First, create the MySQL database and tables by running the provided schema:

```bash
mysql -u root -p < schema.sql
```

### 2. Dependencies

The MySQL driver is already included in the project dependencies:

```go
github.com/go-sql-driver/mysql v1.9.3
```

## Usage

### Basic Usage

```go
package main

import (
    "log"
    
    "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
    "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/drivers/mysql"
)

func main() {
    // Configure database connection
    config := mysql.DatabaseConfig{
        Host:     "localhost",
        Port:     "3306",
        Username: "root",
        Password: "password",
        Database: "ecommerce",
    }

    // Create database connection
    db, err := mysql.NewConnection(config)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create repository
    repository := mysql.NewMySQLBasketRepository(db)

    // Use the repository
    basketFactory := entities.NewBasketFactory()
    basket, err := basketFactory.NewBasket("user123")
    if err != nil {
        log.Fatal(err)
    }

    // Add items
    basket.AddItem("product1", 2)
    basket.AddItem("product2", 1)

    // Save basket
    basketID, err := repository.Save(basket)
    if err != nil {
        log.Fatal(err)
    }

    // Find basket
    foundBasket, err := repository.Find(basketID)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Found basket for user: %s", foundBasket.GetUserID())
}
```

### Using DSN

```go
// Alternative: use DSN string
dsn := "root:password@tcp(localhost:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
db, err := mysql.NewConnectionFromDSN(dsn)
if err != nil {
    log.Fatal(err)
}
```

## Database Schema

The adapter uses two main tables:

### `baskets` table
- `id` (VARCHAR(36), PRIMARY KEY): Unique basket identifier
- `userid` (VARCHAR(255)): User identifier
- `items` (JSON): Serialized basket items
- `created_at` (TIMESTAMP): Creation timestamp
- `updated_at` (TIMESTAMP): Last update timestamp

### Items Storage
Basket items are stored as JSON in the `items` column with the following structure:

```json
{
  "product1": {"product_id": "product1", "count": 2},
  "product2": {"product_id": "product2", "count": 1}
}
```

## Features

- **Upsert Support**: The `Save` method uses `INSERT ... ON DUPLICATE KEY UPDATE` for efficient upsert operations
- **Connection Pooling**: Configured with appropriate connection pool settings
- **Error Handling**: Proper error handling with meaningful error messages
- **JSON Storage**: Efficient storage of basket items using MySQL's JSON data type
- **Indexing**: Optimized with indexes on frequently queried fields

## Testing

Run the example test to verify the setup:

```bash
go test ./internal/domain/basket/drivers/mysql/... -v
```

Note: The test requires a running MySQL database and will be skipped if not available.