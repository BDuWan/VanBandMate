package models

type District struct {
	Code                 string `gorm:"primaryKey;size:20" json:"code"`
	Name                 string `gorm:"size:255" json:"name"`
	NameEn               string `gorm:"size:255" json:"name_en"`
	FullName             string `gorm:"size:255" json:"full_name"`
	FullNameEn           string `gorm:"size:255" json:"full_name_en"`
	CodeName             string `gorm:"size:255" json:"code_name"`
	ProvinceCode         string `gorm:"size:20" json:"province_code"`
	AdministrativeUnitID int    `json:"administrative_unit_id"`
}
