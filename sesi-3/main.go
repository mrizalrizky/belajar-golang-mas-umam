package main

import (
	"fmt"
	"log"
	"mrizalrizky/sesi-3/database"
	"mrizalrizky/sesi-3/handlers"
	"mrizalrizky/sesi-3/models"
	"mrizalrizky/sesi-3/repositories"
	"mrizalrizky/sesi-3/services"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var categories = []models.Category{
	{ID: 1, Name: "Furnitures", Description: "This is furnitures description" },
	{ID: 2, Name: "Smart Home Devices", Description: "This is smart home devices description" },
}

func checkUniqueName(name string) bool {
	for _, category := range categories {
		if category.Name == name {
			return false
		}
	}

	return true
}

type Config struct {
	Port string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		v.SetConfigFile(".env")
		if err := v.ReadInConfig(); err != nil {
			log.Fatal("Failed to read config file:", err)
		}
	}

	config := Config{
		Port: v.GetString("PORT"),
		DBConn: v.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	router := http.NewServeMux()

	
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	router.HandleFunc("/api/v1/categories", categoryHandler.HandleCategories)
	router.HandleFunc("/api/v1/categories/{id}", categoryHandler.HandleCategoryByID)
	
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	router.HandleFunc("/api/v1/products", productHandler.HandleProducts)
	router.HandleFunc("/api/v1/products/{id}", productHandler.HandleProductByID)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	router.HandleFunc("/api/v1/checkout", transactionHandler.HandleCheckout) // POST

	server := http.Server{
		Addr: ":"+config.Port,
		Handler: router,

	} 

	fmt.Println("Server is running on port " + config.Port)
	server.ListenAndServe()
	
}