package handlers

import (
	"encoding/json"
	"mrizalrizky/sesi-3/internal/response"
	"mrizalrizky/sesi-3/models"
	"mrizalrizky/sesi-3/services"
	"net/http"
)


type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}


func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			h.Checkout(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(response.ApiResponse{
				Success: false,
				Message: "Method not allowed",
			})
	}
}

// POST /api/v1/checkout
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	var newCheckout models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&newCheckout); err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Errors: err,
		})
		return
	}

	transaction, err := h.service.Checkout(newCheckout.Items, false)
	if err != nil {
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Error creating checkout",
			Errors: err.Error(),
		})
		return
	}

	encoder.Encode(response.ApiResponse{
		Success: true,
		Message: "Checkout created successfully",
		Data: transaction,
	})
}