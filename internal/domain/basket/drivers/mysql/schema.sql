-- MySQL schema for basket repository
-- Run this script to create the necessary tables

CREATE DATABASE IF NOT EXISTS ecommerce;
USE ecommerce;

-- Baskets table
CREATE TABLE IF NOT EXISTS baskets (
    id VARCHAR(36) PRIMARY KEY,
    userid VARCHAR(255) NOT NULL,
    items JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_userid (userid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Example of the JSON structure stored in the items column:
-- {
--   "product1": {"product_id": "product1", "count": 2},
--   "product2": {"product_id": "product2", "count": 1}
-- }