package handlers

import (
	"encoding/json"
	"mrizalrizky/sesi-3/internal/response"
	"mrizalrizky/sesi-3/services"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GET /api/report/hari-ini
func (h *ReportHandler) HandleTodayReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	report, err := h.service.GetTodaySalesReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Error getting today's report",
			Errors:  err.Error(),
		})
		return
	}

	encoder.Encode(report)
}

// GET /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" && endDateStr == "" {
		report, err := h.service.GetTodaySalesReport()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(response.ApiResponse{
				Success: false,
				Message: "Error getting report",
				Errors:  err.Error(),
			})
			return
		}
		encoder.Encode(report)
		return
	}

	if startDateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "start_date is required when using date range",
		})
		return
	}

	if endDateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "end_date is required when using date range",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid start_date format. Use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Invalid end_date format. Use YYYY-MM-DD",
		})
		return
	}

	if startDate.After(endDate) {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "start_date must be before or equal to end_date",
		})
		return
	}

	report, err := h.service.GetSalesReportByDateRange(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(response.ApiResponse{
			Success: false,
			Message: "Error getting report",
			Errors:  err.Error(),
		})
		return
	}

	encoder.Encode(report)
}
