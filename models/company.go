package models

type Company struct {
	Id                string
	UserId            string `json:"user_id" gorm:"unique"`
	SupportEmail      string `json:"support_email"`
	Name              string `json:"name"  gorm:"unique;default:null" binding:"required"`
	Description       string `json:"description,omitempty" binding:"max=3000"`
	AmountOfEmployees int    `json:"amount_of_employees"  gorm:"unique" binding:"required"`
	IsRegistered      bool   `json:"is_registered" binding:"required"`
	Type              string `json:"type" binding:"required"`
}
