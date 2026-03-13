package service

import (
	"context"

	"pos-system/backend/internal/domain"
	"pos-system/backend/internal/repository"
)

type ReportService struct {
	reports           *repository.ReportRepository
	lowStockThreshold int
}

func NewReportService(reports *repository.ReportRepository, lowStockThreshold int) *ReportService {
	return &ReportService{
		reports:           reports,
		lowStockThreshold: lowStockThreshold,
	}
}

func (s *ReportService) Dashboard(ctx context.Context) (domain.DashboardMetrics, error) {
	return s.reports.GetDashboard(ctx, s.lowStockThreshold)
}

func (s *ReportService) TopProducts(ctx context.Context, dateFrom, dateTo string, limit int) ([]domain.TopProduct, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	return s.reports.GetTopProducts(ctx, dateFrom, dateTo, limit)
}

func (s *ReportService) LowStock(ctx context.Context, threshold int) ([]domain.Product, error) {
	if threshold <= 0 {
		threshold = s.lowStockThreshold
	}
	return s.reports.GetLowStockProducts(ctx, threshold)
}
