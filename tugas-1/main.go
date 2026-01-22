package main

import (
	"encoding/json"
	"fmt"
	"mrizalrizky/tugas-1/internal/model"
	"mrizalrizky/tugas-1/internal/response"
	"net/http"
	"strconv"
)

var categories = []model.Category{
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

func getAllCategories(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if len(categories) <= 0 {
		encoder.Encode(response.ApiResponse{
			Success: true,
			Message: "There are no categories found",
			Data: categories,
		})
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Categories retrieved successfully",
		Data: categories,
	})
}

func getCategoryById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Category ID not found",
			Data:    nil,
		})
		return
	}

	for _, category := range categories {
		if category.ID == id {
			encoder.Encode(response.ApiResponse{
				Success: true,
				Message: "Category retrieved successfully",
				Data:    category,
			})
			return
		}
	}

	encoder.Encode(response.ApiResponse{
		Success: false,
		Message: "Category with that ID not found",
		Data:    nil,
	})
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	var newCategory model.Category

	// err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Error creating category",
		})
		return
	}
	
	if isUnique := checkUniqueName(newCategory.Name); !isUnique {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Category name already exists",
		})
		return
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)
	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Category created successfully",
		Data: newCategory,
	})
}

func updateCategoryById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Category ID not found",
			Data:    nil,
		})
		return
	}

	var updateCategory model.Category
	if err := json.NewDecoder(r.Body).Decode(&updateCategory); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Error updating category",
		})
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			updateCategory.ID = id
			categories[i] = updateCategory
			encoder.Encode(response.ApiResponse{
				Success: true,
				Message: "Category updated successfully",
				Data: updateCategory,
			})
			return
		}
	}

	encoder.Encode(response.ApiResponse{
		Success: false,
		Message: "Category with that ID not found",
	})
}

func deleteCategoryById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Category ID not found",
			Data:    nil,
		})
		return
	}
	for i := range categories {
		if categories[i].ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			encoder.Encode(response.ApiResponse{
				Success: true,
				Message: "Category deleted successfully",
			})
			return
		}
	}

	encoder.Encode(response.ApiResponse{
		Success: false,
		Message: "Category with that ID not found",
	})
}

func main() {
	router := http.NewServeMux()
	
	router.HandleFunc("GET /categories", getAllCategories)
	router.HandleFunc("GET /categories/{id}", getCategoryById)
	router.HandleFunc("POST /categories", createCategory)
	router.HandleFunc("PUT /categories/{id}", updateCategoryById)
	router.HandleFunc("DELETE /categories/{id}", deleteCategoryById)

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	fmt.Println("Server is running on port 8080")
	server.ListenAndServe()
	
}