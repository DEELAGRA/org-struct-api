package models

import "time"

type Departament struct {
	ID         int
	Name       string
	Parent_id  int
	Created_at time.Time
}

type Employee struct {
	ID             int
	Departament_id *Departament
	Full_name      string
	Position       string
	Hired_at       time.Time
}
