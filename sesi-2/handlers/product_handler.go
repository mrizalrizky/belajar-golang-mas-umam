package handlers

import (
	"encoding/json"
	"fmt"
	"mrizalrizky/sesi-2/internal/response"
	"mrizalrizky/sesi-2/models"
	"mrizalrizky/sesi-2/services"
	"net/http"
	"strconv"
)

/*
	Konsep Dependency Injection
	-> Di sebuah resto ada yang mau beli, dan customer memanggil pelayan
	-> Pelayan = Handler
	-> Pelayan ngomong ke koki buat masakin makanan
	-> Koki = Service
	-> Koki minta bahan makanan ke bagian gudang
	-> Bagian Gudang = Repository / Database
*/


type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}



func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
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

// GET /api/products
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, _ := h.service.GetAllProducts()

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if len(products) <= 0 {
		encoder.Encode(response.ApiResponse{
			Success: true,
			Message: "There are no products found",
			Data: products,
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Products retrieved successfully",
		Data: products,
	})
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	var newProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Errors: err,
		})
		return
	}

	if err := h.service.CreateProduct(&newProduct); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Error creating product",
			Errors: err,
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Product created successfully",
		Data: newProduct,
	})
}

// GET/PUT/DELETE /api/v1/products/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
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

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
	
	Product, err := h.service.GetProductByID(id)
	if Product == nil {
		encoder.Encode(response.ApiResponse{
			Success: true,
			Message: "Product with that ID not found",
			Data: nil,
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Data: Product,
	})
}

func (h *ProductHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
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

	var updateProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&updateProduct); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	updateProduct.ID = id
	fmt.Println("UPODATEPRODUCT", updateProduct)
	if err := h.service.UpdateProductByID(&updateProduct); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Failed to update product",
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Data: updateProduct,
	})
}

func (h *ProductHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Product ID not found",
		})
		return
	}

	if err := h.service.DeleteProductByID(id); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Failed to update product",
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}