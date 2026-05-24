package models

import (
	"time"
)

type Department struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;type:varchar(200);not null"`
	ParentID  *int      `gorm:"column:parent_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	Parent    *Department  `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Children  []Department `gorm:"foreignKey:ParentID"`
	Employees []Employee   `gorm:"foreignKey:DepartmentID"`
}

type Employee struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement"`
	DepartmentID int        `gorm:"column:department_id;not null"`
	FullName     string     `gorm:"column:full_name;type:varchar(200);not null"`
	Position     string     `gorm:"column:position;type:varchar(200);not null"`
	HiredAt      *time.Time `gorm:"column:hired_at;type:date"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`

	Department Department `gorm:"foreignKey:DepartmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
