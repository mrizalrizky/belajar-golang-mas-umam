package services

import (
	"mrizalrizky/sesi-3/models"
	"mrizalrizky/sesi-3/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetTodaySalesReport() (*models.SalesReport, error) {
	return s.repo.GetTodaySalesReport()
}

func (s *ReportService) GetSalesReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	return s.repo.GetSalesReportByDateRange(startDate, endDate)
}
