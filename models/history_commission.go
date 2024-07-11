package models

import "time"

type HistoryCommission struct {
	HistoryCommissionID int       `json:"history_commission_id"`
	UserID              int       `json:"user_id"`
	CommissionTotal     float64   `json:"commission_total"`
	CommissionPaid      float64   `json:"commission_paid"`
	CommissionDebt      float64   `json:"commission_debt"`
	Type                string    `json:"type"`
	Description         string    `json:"description"`
	Deleted             bool      `json:"deleted"`
	DeletedBy           int       `json:"deleted_by"`
	UpdatedBy           int       `json:"updated_by"`
	CreatedBy           int       `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt           time.Time `json:"deleted_at"`
}
