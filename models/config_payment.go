package models

import "time"

type ConfigPayment struct {
	ConfigPaymentID   int       `json:"config_payment_id" gorm:"primaryKey;autoIncrement"`
	CommissionDefault float64   `json:"commission_default"`
	CommissionBonus   float64   `json:"commission_bonus"`
	NumberToGetBonus  int       `json:"number_to_get_bonus"`
	NumberDayPeriod   int       `json:"number_day_period"`
	UpdatedAt         time.Time `jaon:"updated_at"`
}
