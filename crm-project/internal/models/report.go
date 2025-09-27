// Create new file: internal/models/report.go
package models

// EmployeeLeadReport is the top-level structure for our JSON response.
type EmployeeLeadReport struct {
	Rows  []EmployeeLeadRow `json:"rows"`
	Total LeadStatusSummary `json:"total"`
}

// EmployeeLeadRow represents the data for a single employee in the report.
type EmployeeLeadRow struct {
	EmployeeID   int               `json:"employee_id"`
	EmployeeName string            `json:"employee_name"`
	Counts       LeadStatusSummary `json:"counts"`
}

// LeadStatusSummary holds the counts for each status. It's used for both individual rows and the total.
type LeadStatusSummary struct {
	New       int `json:"new"`
	Contacted int `json:"contacted"`
	Qualified int `json:"qualified"`
	Converted int `json:"converted"`
	Lost      int `json:"lost"`
}

// DealsPipelineReport represents the structure for the deals pipeline report.
type DealsPipelineReport struct {
	Rows  []DealsPipelineReportRow `json:"rows"`
	Total DealsPipelineSummary   `json:"total"`
}

// DealsPipelineReportRow represents a single row in the deals pipeline report.
type DealsPipelineReportRow struct {
	StageName   string  `json:"stage_name"`
	DealCount   int     `json:"deal_count"`
	TotalAmount float64 `json:"total_amount"`
}

// DealsPipelineSummary holds the total counts and amounts for the deals pipeline.
type DealsPipelineSummary struct {
	TotalDealCount   int     `json:"total_deal_count"`
	TotalDealAmount float64 `json:"total_deal_amount"`
}