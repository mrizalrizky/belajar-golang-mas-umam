// Go -> turunan dari C
// Concurrency -> menjalan sebuah task secara bergantian dengan cepat
// Goroutine -> unit kecil eksekusi yang dibikin oleh Go untuk menjalankan tasknya.
// Ketika ada task yang masuk ke server, Go akan menjalankan Goroutinenya di atas thread

package main

// "net/http" // -> fundamentalnya utk bikin API, bawaan go dan bukan framework

import (
	"encoding/json"
	"fmt"
	"mrizalrizky/belajar-go/internal/model"
	"mrizalrizky/belajar-go/internal/validation"
	"mrizalrizky/belajar-go/pkg/response"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var products = []model.Product{
	// {ID: 1, Name: "Produk 1", Price: 10000, Stock: 21},
	// {ID: 2, Name: "Produk 2", Price: 15000, Stock: 5},
}

func createProduct(product model.Product) (model.Product) {
	// Bac data dari request
	// Masukin data ke dalam variable product
	// var createdProduct model.Product
	
	// err := json.NewDecoder(r.Body).Decode(&createdProduct)
	// if err != nil {
	// 	return model.Product{}
	// 	// http.Error(w, "Error to create product", http.StatusBadRequest)
	// }

	product.ID = len(products) + 1
	products = append(products, product)
	return product
}

func getAllProducts() (*[]model.Product) {
	if(len(products) == 0) {
		return nil
	}
	return &products
}

func getProductById(productId int) (*model.Product) {
	for _, p := range products {
		if p.ID == productId {
			return &p
		}
	}

	return nil
}

func updateProductById(r *http.Request, productId int) (*model.Product) {
	var updatedProduct model.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		return &model.Product{}
	}

	for i := range products {
		if (products[i].ID == productId) {
			updatedProduct.ID = productId
			products[i] = updatedProduct
			
			return &products[i]
		}
	}

	return nil
}

func deleteProductById(r *http.Request, productId int) (*model.Product) {
	for i, p := range products {
		if (p.ID == productId) {
			products = append(products[:i], products[i+1:]...) // [:i] => index 0 sampe i

			return &p
		}
	}

	return nil
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("OK"))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"success": "true",
		"message": "API is running and healthy",
	})
}

// Main -> akan dijalankan pertama kali
func main() {
	router := http.NewServeMux()
	// GET localhost:8080/health
	router.HandleFunc("/api/health", checkHealth)

	// GET localhost:8080/api/products
	// POST localhost:8080/api/products
	router.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
			case "GET":
				data := getAllProducts()
				if(data == nil) {
					encoder.Encode(response.ApiResponse{
						Success: false,
						Message: "No products available",
						Data:    nil,
					})
					return
				}
				encoder.Encode(response.ApiResponse{
					Success: true,
					Message: "Products retrieved successfully",
					Data:    data,
				})

			case "POST":
				var product model.Product
				if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
					http.Error(w, "Error to create product", http.StatusBadRequest)
					return
				}

				err := validate.Struct(product)
				if err != nil {
					encoder.Encode(response.ApiResponse{
						Success: false,
						Message: "Validation error",
						Errors: validation.FormatError(err),
					})
					return
				}
				data := createProduct(product)
				encoder.Encode(response.ApiResponse{
					Success: true,
					Message: "Product created successfully",
					Data:    data,
				})
		}
	})

	// // GET localhost:8080/api/products/:id
	// // PUT/PATCH localhost:8080/api/products/:id
	// // DELETE localhost:8080/api/products/:id
	router.HandleFunc("/api/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			// http.Error(w, "Product ID not found", http.StatusNotFound)
			encoder.Encode(response.ApiResponse{
				Success: false,
				Message: "Product ID not found",
				Data:    nil,
			})
			return
		}

		switch r.Method {
			case "GET":
				data := getProductById(id)
				if(data == nil) {
					encoder.Encode(response.ApiResponse{
						Success: false,
						Message: "Product with that ID not found",
						Data:    nil,
					})
					return
				}
				encoder.Encode(response.ApiResponse{
					Success: true,
					Message: "Products retrieved successfully",
					Data:    data,	
				})
				return
			case "PUT":
				valid := validate.Struct(r.Body)
				fmt.Println("VALIDDD", valid)
				data := updateProductById(r, id)
				if(data == nil) {
					encoder.Encode(response.ApiResponse{
						Success: false,
						Message: "Product with that ID not found",
						Data:    nil,
					})
					return
				}
				encoder.Encode(response.ApiResponse{
					Success: true,
					Message: "Product updated successfully",
					Data:    data,	
				})
				return
			case "DELETE":
				data := deleteProductById(r, id)
				if(data == nil) {
					encoder.Encode(response.ApiResponse{
						Success: false,
						Message: "Product with that ID not found",
						Data:    nil,
					})
					return
				}
				encoder.Encode(response.ApiResponse{
					Success: true,
					Message: "Product deleted successfully",
					Data:    data,	
				})
				return	
		}

	})

	
	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	}
	
	fmt.Println("Server is running on port 8080")
	server.ListenAndServe()
}

