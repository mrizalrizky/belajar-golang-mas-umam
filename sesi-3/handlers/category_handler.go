package handlers

import (
	"encoding/json"
	"mrizalrizky/sesi-3/internal/response"
	"mrizalrizky/sesi-3/models"
	"mrizalrizky/sesi-3/services"
	"net/http"
	"strconv"
)


type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}


func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			h.GetAll(w, r)

		case http.MethodPost:
			h.Create(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(response.ApiResponse{
				Success: false,
				Message: "Method not allowed",
			})
	}
}

// GET /api/categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, _ := h.service.GetAllCategories()

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if len(categories) <= 0 {
		encoder.Encode(response.ApiResponse{
			Success: true,
			Message: "There are no categories found",
			Data: categories,
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Categories retrieved successfully",
		Data: categories,
	})
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	var newCategory models.Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Errors: err,
		})
		return
	}

	if err := h.service.CreateCategory(&newCategory); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Error creating category",
			Errors: err,
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Category created successfully",
		Data: newCategory,
	})
}

// GET/PUT/DELETE /api/v1/categories/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			h.GetByID(w, r)

		case http.MethodPut:
			h.UpdateByID(w, r)
		
		case http.MethodDelete:
			h.DeleteByID(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(response.ApiResponse{
				Success: false,
				Message: "Method not allowed",
			})
	}
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}
	
	category, err := h.service.GetCategoryByID(id)
	if category == nil {
		encoder.Encode(response.ApiResponse{
			Success: true,
			Message: "Category with that ID not found",
			Data: nil,
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Category retrieved successfully",
		Data: category,
	})
}

func (h *CategoryHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	var updateCategory models.Category
	if err := json.NewDecoder(r.Body).Decode(&updateCategory); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	updateCategory.ID = id
	if err := h.service.UpdateCategoryByID(&updateCategory); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Failed to update category",
			Errors: err,
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Category retrieved successfully",
		Data: updateCategory,
	})
}

func (h *CategoryHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.DeleteCategoryByID(id); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Failed to update category",
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}