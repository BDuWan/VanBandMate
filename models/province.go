package models

type Province struct {
	Code                   string `gorm:"primaryKey;size:20" json:"code"`
	Name                   string `gorm:"size:255" json:"name"`
	NameEn                 string `gorm:"size:255" json:"name_en"`
	FullName               string `gorm:"size:255" json:"full_name"`
	FullNameEn             string `gorm:"size:255" json:"full_name_en"`
	CodeName               string `gorm:"size:255" json:"code_name"`
	AdministrativeUnitID   int    `json:"administrative_unit_id"`
	AdministrativeRegionID int    `json:"administrative_region_id"`
}
