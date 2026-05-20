package models

import (
	"time"
)

type Department struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Parent_id  int       `gorm:"index" json:"parent_id"`
	Created_at time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Department) TableName() string {
	return "department"
}

type Employee struct {
	ID             int         `json:"id"`
	Departament_id *Department `json:"departament_id"`
	Full_name      string      `json:"full_name"`
	Position       string      `json:"position"`
	Hired_at       time.Time   `json:"hired_at"`
}

func (Employee) TableName() string {
	return "employee"
}
