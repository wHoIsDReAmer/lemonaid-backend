package db

import (
	"gorm.io/gorm"
	"time"
)

// PLAN
const (
	STANDARD   = 1
	PREMIUM    = 2
	RESUME     = 3
	SPECIALIST = 4
)

// USER TYPE
const (
	ACADEMY = 1
	TEACHER = 2
	ADMIN   = 3
)

type User struct {
	gorm.Model  `json:"-"`
	FirstName   string    `gorm:"size:255" json:"first_name"`
	LastName    string    `gorm:"size:255" json:"last_name"`
	Email       string    `gorm:"size:320;unique" json:"email"`
	Password    string    `gorm:"size:255" json:"-"`
	Salt        string    `gorm:"size:10" json:"-"`
	PhoneNumber string    `gorm:"size:15;unique" json:"phone_number"`
	Birthday    time.Time `gorm:"type:time" json:"birthday"`
	Gender      *string   `gorm:"size:30" json:"gender"`
	Nationality *string   `gorm:"size:100;default:'Republic of Korea'" json:"nationality"`
	VisaCode    *string   `gorm:"size:16" json:"visa_code"`
	Occupation  *string   `gorm:"size:32" json:"occupation"`

	// Will be deleted
	Manners int `gorm:"default:0" json:"manners"`
	// Will be deleted
	Amateur int `gorm:"default:0" json:"amateur"`

	VideoMessenger   *string `gorm:"size:255" json:"video_messenger"`
	VideoMessengerID *string `gorm:"size:255" json:"video_messenger_id"`
	Resume           *[]byte `gorm:"type:longblob" json:"resume"`
	Image            *string `gorm:"size:255" json:"image_path"`
	Plan             int     `gorm:"default:0" json:"plan"`

	UserAccepted int `gorm:"size:1;default:0" json:"-"`
	UserType     int `gorm:"default:0" json:"user_type"`
}

type Tour struct {
	gorm.Model

	// Seperate by ","
	Images *string `gorm:"size:255" json:"images"`

	TourName    string `gorm:"size:255" json:"tour_name"`
	Description string `gorm:"size:255" json:"description"`
	PostOwn     string `gorm:"size:255" json:"post_own"`
	Company     string `gorm:"size:255" json:"company"`
	Theme       string `gorm:"size:255" json:"theme"`
	Location    string `gorm:"size:255" json:"location"`
	Date        string `gorm:"type:date" json:"date"`
	Price       uint   `json:"price"`
	Itinerary   string `gorm:"type:text" json:"itinerary"`
}

type JobPost struct {
	gorm.Model

	Academy  string `gorm:"size:255" json:"academy"`
	Campus   string `gorm:"size:255" json:"campus"`
	Category string `gorm:"size:255" json:"category"`

	// Seperate by ","
	Images *string `gorm:"size:255" json:"images"`

	Location     string `gorm:"size255" json:"location"`
	Position     string `gorm:"size:255" json:"position"`
	SalaryMin    string `json:"start_salary"`
	SalaryMax    string `json:"end_salary"`
	StudentLevel string `gorm:"size:255" json:"student_level"`

	// Two column types will be to string
	WorkingHoursMin string `json:"working_hours_start"`
	WorkingHoursMax string `json:"working_hours_end"`

	PaidVacation     uint   `json:"paid_vacation"`
	AnnualLeave      uint   `json:"annual_leave"`
	Severance        string `gorm:"size:255" json:"severance"`
	Insurance        string `gorm:"size:255" json:"insurance"`
	Housing          string `gorm:"size:255" json:"housing"`
	HousingAllowance string `gorm:"size:255" json:"housing_allowance"`
	Rank             int    `gorm:"size:1" json:"rank"`

	UserID uint
	User   User `gorm:"ForeignKey:UserID;References:ID"`
}

type PendingJobPost struct {
	gorm.Model

	Academy  string `gorm:"size:255" json:"academy"`
	Campus   string `gorm:"size:255" json:"campus"`
	Category string `gorm:"size:255" json:"category"`

	// Seperate by ","
	Images *string `gorm:"size:255" json:"images"`

	Location     string `gorm:"size255" json:"location"`
	Position     string `gorm:"size:255" json:"position"`
	SalaryMin    string `json:"start_salary"`
	SalaryMax    string `json:"end_salary"`
	StudentLevel string `gorm:"size:255" json:"student_level"`

	// Two column types will be to string
	WorkingHoursMin string `json:"working_hours_start"`
	WorkingHoursMax string `json:"working_hours_end"`

	PaidVacation     uint   `json:"paid_vacation"`
	AnnualLeave      uint   `json:"annual_leave"`
	Severance        string `gorm:"size:255" json:"severance"`
	Insurance        string `gorm:"size:255" json:"insurance"`
	Housing          string `gorm:"size:255" json:"housing"`
	HousingAllowance string `gorm:"size:255" json:"housing_allowance"`
	Rank             int    `gorm:"size:1" json:"rank"`

	UserID uint
	User   User `gorm:"ForeignKey:UserID;References:ID"`
}

type PartyAndEvents struct {
	gorm.Model

	// Seperate by ","
	Images *string `gorm:"size:255" json:"images"`

	PartyName   string `gorm:"size:255" json:"party_name"`
	Description string `gorm:"size:255" json:"description"`
	PostOwn     string `gorm:"size:255" json:"post_own"`
	Company     string `gorm:"size:255" json:"company"`
	Theme       string `gorm:"size:255" json:"theme"`
	Location    string `gorm:"size:255" json:"location"`
	Date        string `gorm:"type:date" json:"date"`
	Price       uint   `json:"price"`
	Itinerary   string `gorm:"type:text" json:"itinerary"`
}

type Session struct {
	gorm.Model
	Email   string    `gorm:"size:255;unique"`
	Uuid    string    `gorm:"size:36"`
	Expires time.Time `gorm:"type:time"`
}

type ApplyJobPost struct {
	gorm.Model
	JobPostID uint
	JobPost   JobPost `gorm:"ForeignKey:JobPostID;References:ID"`
	UserID    uint
	User      User `gorm:"ForeignKey:UserID;References:ID"`
}
