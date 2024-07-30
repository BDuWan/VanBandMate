package structs

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
	Image             string `form:"image"`
	RoleID            int    `form:"role_id"`
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

type AccUser struct {
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	Email                  string `json:"email"`
	PhoneNumber            string `json:"phone_number"`
	Address                string `json:"address"`
	Password               string `json:"password"`
	ReferralCode           string `json:"referral_code"`
	NameBusiness           string `json:"name_business"`
	FullNameRepresentative string `json:"full_name_representative"`
	TypeUserID             int    `json:"type_user_id"`
	RoleID                 int    `json:"role_id"`
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
	Name         string `form:"name" validate:"required,min=3,max=30"`
	PermissionID []int  `form:"permission_id" validate:"required"`
}

type FormStateUser struct {
	ID    int  `json:"id"`
	State bool `json:"state"`
}

// type Course struct {
// 	Name        string `json:"name"`
// 	CourseID    int    `json:"course_id"`
// 	UserID      int    `json:"user_id"`
// 	Description string `json:"description"`
// }

//	type Course struct {
//		Name           string `json:"name"`
//		StudyProgramID int    `json:"study_program_id"`
//		UserID         int    `json:"user_id"`
//		Price          string `json:"price"`
//		Description    string `json:"description"`
//	}
type Course struct {
	Name           string `json:"name"`
	StudyProgramID int    `json:"study_program_id"`
	UserID         int    `json:"user_id"`
	Description    string `json:"description"`
}
type CourseInstructor struct {
	CourseID int `json:"course_id"`
	UserID   int `json:"user_id"`
}

type StudyProgramInstructor struct {
	StudyProgramID int `json:"study_program_id"`
	UserID         int `json:"user_id"`
}

type FromChangePass struct {
	CurrentPass string `json:"current_pass"`
	NewPass     string `json:"new_pass"`
}

// type CourseUser struct {
// 	CourseID int `json:"course_id"`
// 	UserID   int `json:"user_id"`
// }

type StudyProgramUser struct {
	StudyProgramID int `json:"study_program_id"`
	UserID         int `json:"user_id"`
}
type CourseUser struct {
	CourseID int `json:"course_id"`
	UserID   int `json:"user_id"`
}

type Lesson struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CourseID    int    `json:"course_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	LinkRecord  string `json:"link_record"`
	LinkStudy   string `json:"link_study"`
}

type Lesson1 struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CourseID    int    `json:"course_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	LinkRecord  string `json:"link_record"`
	LinkStudy   string `json:"link_study"`
}
type Assignment struct {
	Name        string `json:"name"`
	CourseID    int    `json:"course_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
	//StartTime time.Time `json:"start_time"`
	//EndTime   time.Time `json:"end_time"`
}
type Assignment1 struct {
	Name        string `json:"name"`
	CourseID    int    `json:"clas_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
	//StartTime time.Time `json:"start_time"`
	//EndTime   time.Time `json:"end_time"`
}
type AssignmentCourse struct {
	CourseID     int    `json:"course_id"`
	UserID       int    `json:"user_id"`
	AssignmentID int    `json:"assignment_id"`
	Status       int    `json:"status"`
	Result       string `json:"result"`
}

type ConfigPayment struct {
	CommissionDefault float64 `json:"commission_default"`
	CommissionBonus   float64 `json:"commission_bonus"`
	NumberStudent     int     `json:"number_student"`
	NumberDay         int     `json:"number_day"`
}

type PriceProgram struct {
	Price       int    `json:"price"`
	Commission  int    `json:"commission"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type PriceProgramUser struct {
	PriceUserID            int    `json:"price_user_id"`
	PriceProgramID         int    `json:"price_program_id"`
	Price                  int    `json:"price"`
	Commission             int    `json:"commission"`
	UserID                 int    `json:"user_id"`
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	NameBusiness           string `json:"name_business"`
	FullNameRepresentative string `json:"full_name_representative"`
}

type Payment struct {
	PriceProgramID int `json:"price_program_id"`
	UserID         int `json:"user_id"`
	Total          int `json:"total"`
}

type CommissionEdit struct {
	UserID     int `json:"user_id"`
	Commission int `json:"commission"`
}

type CommissionPay struct {
	UserID int     `json:"user_id"`
	Price  float64 `json:"price"`
}

//type DataWS struct {
//	AssignmentN int `json:"assignment_n"`
//	LessonN     int `json:"lesson_n"`
//	Assignments []models.AssignmentUser
//	Lessons     []models.Lesson
//}

//type TopCourseUser struct {
//	CourseID int    `json:"course_id"`
//	Title    string `json:"title"`
//	Number   int    `json:"number"`
//}

type TopSaleUser struct {
	ReferralCode string `json:"referral_code"`
	Name         string `json:"name"`
	Number       int    `json:"number"`
}
