package main

import (
	"fmt"
	"log"
	"mrizalrizky/sesi-2/database"
	"mrizalrizky/sesi-2/handlers"
	"mrizalrizky/sesi-2/models"
	"mrizalrizky/sesi-2/repositories"
	"mrizalrizky/sesi-2/services"
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

// func getAllCategories(w http.ResponseWriter, r *http.Request) {
// 	encoder := json.NewEncoder(w)
// 	w.Header().Set("Content-Type", "application/json")

// 	if len(categories) <= 0 {
// 		encoder.Encode(response.ApiResponse{
// 			Success: true,
// 			Message: "There are no categories found",
// 			Data: categories,
// 		})
// 	}

// 	encoder.Encode(response.ApiResponse{
// 		Success: true,
// 		Message: "Categories retrieved successfully",
// 		Data: categories,
// 	})
// }

// func getCategoryById(w http.ResponseWriter, r *http.Request) {
// 	encoder := json.NewEncoder(w)
// 	w.Header().Set("Content-Type", "application/json")
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		encoder.Encode(response.ApiResponse{
// 			Success: false,
// 			Message: "Category ID not found",
// 			Data:    nil,
// 		})
// 		return
// 	}

// 	for _, category := range categories {
// 		if category.ID == id {
// 			encoder.Encode(response.ApiResponse{
// 				Success: true,
// 				Message: "Category retrieved successfully",
// 				Data:    category,
// 			})
// 			return
// 		}
// 	}

// 	encoder.Encode(response.ApiResponse{
// 		Success: false,
// 		Message: "Category with that ID not found",
// 		Data:    nil,
// 	})
// }

// func createCategory(w http.ResponseWriter, r *http.Request) {
// 	encoder := json.NewEncoder(w)
// 	w.Header().Set("Content-Type", "application/json")

// 	var newCategory models.Category

// 	// err := json.NewDecoder(r.Body).Decode(&newCategory)
// 	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
// 		encoder.Encode(response.ApiResponse{
// 			Success: false,
// 			Message: "Error creating category",
// 		})
// 		return
// 	}
	
// 	if isUnique := checkUniqueName(newCategory.Name); !isUnique {
// 		encoder.Encode(response.ApiResponse{
// 			Success: false,
// 			Message: "Category name already exists",
// 		})
// 		return
// 	}

// 	newCategory.ID = len(categories) + 1
// 	categories = append(categories, newCategory)
// 	encoder.Encode(response.ApiResponse{
// 		Success: true,
// 		Message: "Category created successfully",
// 		Data: newCategory,
// 	})
// }

// func updateCategoryById(w http.ResponseWriter, r *http.Request) {
// 	encoder := json.NewEncoder(w)
// 	w.Header().Set("Content-Type", "application/json")
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		encoder.Encode(response.ApiResponse{
// 			Success: false,
// 			Message: "Category ID not found",
// 			Data:    nil,
// 		})
// 		return
// 	}

// 	var updateCategory models.Category
// 	if err := json.NewDecoder(r.Body).Decode(&updateCategory); err != nil {
// 		encoder.Encode(response.ApiResponse{
// 			Success: false,
// 			Message: "Error updating category",
// 		})
// 		return
// 	}

// 	for i := range categories {
// 		if categories[i].ID == id {
// 			updateCategory.ID = id
// 			categories[i] = updateCategory
// 			encoder.Encode(response.ApiResponse{
// 				Success: true,
// 				Message: "Category updated successfully",
// 				Data: updateCategory,
// 			})
// 			return
// 		}
// 	}

// 	encoder.Encode(response.ApiResponse{
// 		Success: false,
// 		Message: "Category with that ID not found",
// 	})
// }

// func deleteCategoryById(w http.ResponseWriter, r *http.Request) {
// 	encoder := json.NewEncoder(w)
// 	w.Header().Set("Content-Type", "application/json")
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		encoder.Encode(response.ApiResponse{
// 			Success: false,
// 			Message: "Category ID not found",
// 			Data:    nil,
// 		})
// 		return
// 	}
// 	for i := range categories {
// 		if categories[i].ID == id {
// 			categories = append(categories[:i], categories[i+1:]...)
// 			encoder.Encode(response.ApiResponse{
// 				Success: true,
// 				Message: "Category deleted successfully",
// 			})
// 			return
// 		}
// 	}

// 	encoder.Encode(response.ApiResponse{
// 		Success: false,
// 		Message: "Category with that ID not found",
// 	})
// }

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

	server := http.Server{
		Addr: ":"+config.Port,
		Handler: router,

	} 

	fmt.Println("Server is running on port " + config.Port)
	server.ListenAndServe()
	
}