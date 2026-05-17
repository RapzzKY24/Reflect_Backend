package database

import (
	"fmt"
	"log"

	"reflect-backend/internal/config"

	"reflect-backend/internal/modules/address"
	"reflect-backend/internal/modules/cart"
	"reflect-backend/internal/modules/category"
	"reflect-backend/internal/modules/collection"
	"reflect-backend/internal/modules/order"
	"reflect-backend/internal/modules/product"
	"reflect-backend/internal/modules/user"
	"reflect-backend/internal/modules/wishlist"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = database

	log.Println("Database connected successfully")
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&user.User{},
		&product.Product{},
		&category.Category{},
		&collection.Collection{},
		&cart.CartItem{},
		&wishlist.WishlistItem{},
		&address.Address{},
		&order.Order{},
		&order.OrderItem{},
		&order.OrderTracking{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database migrated successfully")
}