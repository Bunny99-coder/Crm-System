// Replace the contents of internal/service/report_service.go
package service

import (
	"context"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"log/slog"
	"sync"
)

type ReportService struct {
	userRepo *postgres.UserRepo
	leadRepo *postgres.LeadRepo
	dealRepo *postgres.DealRepo
	logger   *slog.Logger
}

func NewReportService(ur *postgres.UserRepo, lr *postgres.LeadRepo, dr *postgres.DealRepo, logger *slog.Logger) *ReportService {
	return &ReportService{
		userRepo: ur,
		leadRepo: lr,
		dealRepo: dr,
		logger:   logger,
	}
}

func (s *ReportService) GenerateEmployeeLeadReport(ctx context.Context) (*models.EmployeeLeadReport, error) {
	s.logger.Info("starting generation of employee lead report")
	agents, err := s.userRepo.GetAllSalesAgents(ctx)
	if err != nil {
		s.logger.Error("failed to get sales agents for report", "error", err)
		return nil, err
	}

	var wg sync.WaitGroup
	resultsChan := make(chan models.EmployeeLeadRow, len(agents))
	report := &models.EmployeeLeadReport{}

	for _, agent := range agents {
		wg.Add(1)
		go func(currentAgent models.User) {
			defer wg.Done()
			
			counts, err := s.leadRepo.GetLeadCountsByUserID(ctx, currentAgent.ID)
			if err != nil {
				s.logger.Error("failed to get lead counts for agent", "agent_id", currentAgent.ID, "error", err)
				return
			}
			resultsChan <- models.EmployeeLeadRow{
				EmployeeID:   currentAgent.ID,
				EmployeeName: currentAgent.Username,
				Counts: models.LeadStatusSummary{
					New:       counts.New,
					Contacted: counts.Contacted,
					Qualified: counts.Qualified,
					Converted: counts.Converted,
					Lost:      counts.Lost,
				},
			}
		}(agent)
	}

	wg.Wait()
	close(resultsChan)

	for row := range resultsChan {
		report.Rows = append(report.Rows, row)
		report.Total.New += row.Counts.New
		report.Total.Contacted += row.Counts.Contacted
		report.Total.Qualified += row.Counts.Qualified
		report.Total.Converted += row.Counts.Converted
		report.Total.Lost += row.Counts.Lost
	}

	if ctx.Err() != nil {
		s.logger.Warn("employee lead report generation cancelled by context", "error", ctx.Err())
		return nil, ctx.Err()
	}

	s.logger.Info("successfully generated employee lead report", "row_count", len(report.Rows))
	return report, nil
}

func (s *ReportService) GetSourceLeadReport(ctx context.Context) ([]postgres.SourceLeadReportRow, error) {
	s.logger.Info("generating source lead report")
	return s.leadRepo.GetSourceLeadReport(ctx)
}

func (s *ReportService) GetEmployeeSalesReport(ctx context.Context) ([]postgres.EmployeeSalesReportRow, error) {
	s.logger.Info("generating employee sales report")
	return s.dealRepo.GetEmployeeSalesReport(ctx)
}



func (s *ReportService) GetSourceSalesReport(ctx context.Context) ([]postgres.SourceSalesReportRow, error) {
	s.logger.Info("generating source sales report")
	return s.dealRepo.GetSourceSalesReport(ctx)
}