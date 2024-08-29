package structs

import "lms/models"

type User struct {
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	Email                  string `json:"email"`
	PhoneNumber            string `json:"phone_number"`
	Address                string `json:"address"`
	Password               string `json:"password"`
	ReferralCode           string `json:"referral_code"`
	NameBusiness           string `json:"name_business"`
	FullNameRepresentative string `json:"full_name_representative"`
	RecaptchaResponse      string `json:"recaptcha_response"`
}

type SignUpForm struct {
	FirstName         string `form:"first_name"`
	LastName          string `form:"last_name"`
	PhoneNumber       string `form:"phone_number"`
	LinkFacebook      string `form:"link_facebook"`
	Image             string `form:"image"`
	RoleID            int    `form:"role_id"`
	Gender            int    `form:"gender"`
	Email             string `form:"email"`
	Password          string `form:"password"`
	ConfirmPassword   string `form:"confirm_password"`
	ProvinceCode      string `form:"province_code"`
	DistrictCode      string `form:"district_code"`
	WardCode          string `form:"ward_code"`
	AddressDetail     string `form:"address_detail"`
	DateOfBirth       string `form:"date_of_birth"`
	RecaptchaResponse string `form:"recaptcha_response"`
}

type UpdateInfoForm struct {
	FirstName     string `form:"first_name"`
	LastName      string `form:"last_name"`
	PhoneNumber   string `form:"phone_number"`
	LinkFacebook  string `form:"link_facebook"`
	Image         string `form:"image"`
	Gender        int    `form:"gender"`
	Email         string `form:"email"`
	ProvinceCode  string `form:"province_code"`
	DistrictCode  string `form:"district_code"`
	WardCode      string `form:"ward_code"`
	AddressDetail string `form:"address_detail"`
	DateOfBirth   string `form:"date_of_birth"`
}

type ReqBody struct {
	Draw   int `json:"draw"`
	Start  int `json:"start"`
	Length int `json:"length"`
	Order  []struct {
		Column int    `json:"column"`
		Dir    string `json:"dir"`
	} `json:"order"`
	Search struct {
		Value string `json:"value"`
	} `json:"search"`
}

type RoleForm struct {
	Name        string `form:"name"`
	Describe    string `form:"describe"`
	Permissions []int  `form:"permissions"`
}

type HiringNewsForm struct {
	Province string `json:"province" form:"province"`
	District string `json:"district" form:"district"`
	Ward     string `json:"ward" form:"ward"`
	Date     string `json:"date" form:"date"`
	Price    string `json:"price" form:"price"`
	Address  string `json:"address" form:"address"`
	Describe string `json:"describe" form:"describe"`
}

//type RoleForm struct {
//	Name        string `json:"name"`
//	Describe    string `json:"describe"`
//	Permissions []int  `json:"permissions"`
//}

type FormStateUser struct {
	ID    int  `json:"id"`
	State bool `json:"state"`
}

type FormChangePass struct {
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
	CfPass  string `json:"cf_pass"`
}

type FormUpdatePass struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	CfPassword string `json:"cf_password"`
	Otp        string `json:"otp"`
}

type FormFilter struct {
	Employer     int `json:"employer"`
	HiringEnough int `json:"hiring_enough"`
	Year         int `json:"year"`
	Month        int `json:"month"`
	TimeCreate   int `json:"time_create"`
	Order        int `json:"order"`
	Page         int `json:"page"`
	ItemsPerPage int `json:"items_per_page"`
}

type FormFilterContract struct {
	Status       int `json:"status"`
	Year         int `json:"year"`
	Month        int `json:"month"`
	TimeCreate   int `json:"time_create"`
	Order        int `json:"order"`
	Page         int `json:"page"`
	ItemsPerPage int `json:"items_per_page"`
}

type FormFilterNews struct {
	HiringEnough int    `json:"hiring_enough"`
	Year         int    `json:"year"`
	Month        int    `json:"month"`
	TimeCreate   int    `json:"time_create"`
	ProvinceCode string `json:"province"`
	DistrictCode string `json:"district"`
	WardCode     string `json:"ward"`
	Order        int    `json:"order"`
	Page         int    `json:"page"`
	ItemsPerPage int    `json:"items_per_page"`
}

type FormFind struct {
	Condition    int               `json:"condition"`
	HiringNews   models.HiringNews `json:"hiringNews"`
	Page         int               `json:"page"`
	ItemsPerPage int               `json:"items_per_page"`
}

//	type SelectedItem struct {
//		UserHiringNewsID int `json:"user_hiring_news_id"`
//		NhacCongID       int `json:"nhaccong_id"`
//	}
//
//	type FormSaveApply struct {
//		HiringNewsID  int            `form:"hiring_news_id"`
//		HiringEnough  bool           `form:"hiring_enough"`
//		SelectedItems []SelectedItem `form:"selected_items"`
//	}
type FormSaveApply struct {
	HiringNewsID  int  `json:"hiring_news_id"`
	HiringEnough  bool `json:"hiring_enough"`
	SelectedItems []struct {
		UserHiringNewsID int `json:"user_hiring_news_id"`
		NhaccongID       int `json:"nhaccong_id"`
	} `json:"selected_items"`
}
