package model

// User struct with GORM model
type User struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Email  string `json:"email" gorm:"unique"`
	Salary int    `json:"salary"`
}
