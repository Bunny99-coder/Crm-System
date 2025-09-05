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