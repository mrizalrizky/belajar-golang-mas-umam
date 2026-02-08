package repositories

import (
	"database/sql"
	"mrizalrizky/sesi-3/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetSalesReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	report := &models.SalesReport{}

	summaryQuery := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_transaksi
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`
	err := r.db.QueryRow(summaryQuery, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	topProductQuery := `
		SELECT
			p.name,
			COALESCE(SUM(td.quantity), 0) as qty_terjual
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`
	
	var produkTerlaris models.ProdukTerlaris
	err = r.db.QueryRow(topProductQuery, startDate, endDate).Scan(&produkTerlaris.Nama, &produkTerlaris.QtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	if err != sql.ErrNoRows {
		report.ProdukTerlaris = &produkTerlaris
	}

	return report, nil
}

func (r *ReportRepository) GetTodaySalesReport() (*models.SalesReport, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)
	
	return r.GetSalesReportByDateRange(startOfDay, endOfDay)
}
