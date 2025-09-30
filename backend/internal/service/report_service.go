package service

import (
	"context"
	"crm-project/internal/config"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"crm-project/internal/util"
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

type ReportService struct {
	userRepo *postgres.UserRepo
	leadRepo *postgres.LeadRepo
	dealRepo *postgres.DealRepo
	cfg      *config.Config // Add config here
	logger   *slog.Logger
}

func NewReportService(ur *postgres.UserRepo, lr *postgres.LeadRepo, dr *postgres.DealRepo, cfg *config.Config, logger *slog.Logger) *ReportService {
	return &ReportService{
		userRepo: ur,
		leadRepo: lr,
		dealRepo: dr,
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *ReportService) GenerateEmployeeLeadReport(ctx context.Context) (*models.EmployeeLeadReport, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for GenerateEmployeeLeadReport", "user_id", claims.UserID, "role_id", claims.RoleID)
		return nil, fmt.Errorf("forbidden: only managers can generate this report")
	}

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
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for GetSourceLeadReport", "user_id", claims.UserID, "role_id", claims.RoleID)
		return nil, fmt.Errorf("forbidden: only managers can generate this report")
	}

	s.logger.Info("generating source lead report")
	return s.leadRepo.GetSourceLeadReport(ctx)
}

func (s *ReportService) GetEmployeeSalesReport(ctx context.Context) ([]postgres.EmployeeSalesReportRow, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for GetEmployeeSalesReport", "user_id", claims.UserID, "role_id", claims.RoleID)
		return nil, fmt.Errorf("forbidden: only managers can generate this report")
	}

	s.logger.Info("generating employee sales report")
	return s.dealRepo.GetEmployeeSalesReport(ctx)
}



func (s *ReportService) GetSourceSalesReport(ctx context.Context) ([]postgres.SourceSalesReportRow, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for GetSourceSalesReport", "user_id", claims.UserID, "role_id", claims.RoleID)
		return nil, fmt.Errorf("forbidden: only managers can generate this report")
	}

	s.logger.Info("generating source sales report")
	return s.dealRepo.GetSourceSalesReport(ctx)
}

func (s *ReportService) GetMySalesReport(ctx context.Context) ([]postgres.EmployeeSalesReportRow, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	s.logger.Info("generating personal sales report for user", "user_id", claims.UserID)
	return s.dealRepo.GetEmployeeSalesReportForUser(ctx, claims.UserID)
}

func (s *ReportService) GetDealsPipelineReport(ctx context.Context) (*models.DealsPipelineReport, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for GetDealsPipelineReport", "user_id", claims.UserID, "role_id", claims.RoleID)
		return nil, fmt.Errorf("forbidden: only managers can generate this report")
	}

	s.logger.Info("generating deals pipeline report")
	// Call the repository method to get the raw data
	rawReport, err := s.dealRepo.GetDealsPipelineReport(ctx)
	if err != nil {
		s.logger.Error("failed to get deals pipeline report from repository", "error", err)
		return nil, err
	}

	// Process raw data into the desired report format
	report := &models.DealsPipelineReport{
		Rows: make([]models.DealsPipelineReportRow, len(rawReport)),
	}
	var totalDealCount int
	var totalDealAmount float64

	for i, row := range rawReport {
		report.Rows[i] = models.DealsPipelineReportRow{
			StageName:   row.StageName,
			DealCount:   row.DealCount,
			TotalAmount: row.TotalAmount,
		}
		totalDealCount += row.DealCount
		totalDealAmount += row.TotalAmount
	}

	report.Total = models.DealsPipelineSummary{
		TotalDealCount:   totalDealCount,
		TotalDealAmount: totalDealAmount,
	}

	s.logger.Info("successfully generated deals pipeline report", "row_count", len(report.Rows))
	return report, nil
}
