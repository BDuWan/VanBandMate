package models

type Ward struct {
	Code                 string `gorm:"primaryKey;type:varchar(20)" json:"code"`
	Name                 string `gorm:"type:varchar(255);not null" json:"name"`
	NameEn               string `gorm:"type:varchar(255)" json:"name_en"`
	FullName             string `gorm:"type:varchar(255)" json:"full_name"`
	FullNameEn           string `gorm:"type:varchar(255)" json:"full_name_en"`
	CodeName             string `gorm:"type:varchar(255)" json:"code_name"`
	DistrictCode         string `gorm:"type:varchar(20)" json:"district_code"`
	AdministrativeUnitID int    `gorm:"type:int" json:"administrative_unit_id"`
}
