// Replace the contents of internal/api/handlers/report_handler.go
package handlers

import (
	"crm-project/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ReportHandler struct {
	service *service.ReportService
	logger  *slog.Logger
}

func NewReportHandler(s *service.ReportService, logger *slog.Logger) *ReportHandler {
	return &ReportHandler{service: s, logger: logger}
}

func (h *ReportHandler) GetEmployeeLeadReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report, err := h.service.GenerateEmployeeLeadReport(ctx)
	if err != nil {
		// The service already logged the specific error
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetSourceLeadReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report, err := h.service.GetSourceLeadReport(ctx)
	if err != nil {
		h.logger.Error("failed to generate source lead report", "error", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetEmployeeSalesReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report, err := h.service.GetEmployeeSalesReport(ctx)
	if err != nil {
		h.logger.Error("failed to generate employee sales report", "error", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetSourceSalesReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report, err := h.service.GetSourceSalesReport(ctx)
	if err != nil {
		h.logger.Error("failed to generate source sales report", "error", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetMySalesReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report, err := h.service.GetMySalesReport(ctx)
	if err != nil {
		h.logger.Error("failed to generate personal sales report", "error", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetDealsPipelineReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report, err := h.service.GetDealsPipelineReport(ctx)
	if err != nil {
		h.logger.Error("failed to generate deals pipeline report", "error", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}