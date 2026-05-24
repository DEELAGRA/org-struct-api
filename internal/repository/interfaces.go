package repository

import (
	"context"

	"github.com/DEELAGRA/org-struct-api/models"
)

type DepartmentRepository interface {
	CreateDepartment(ctx context.Context, dept *models.Department) error
	GetDepartment(ctx context.Context, id int) (*models.Department, error)
	UpdateDepartment(ctx context.Context, dept *models.Department) error
	DeleteDepartment(ctx context.Context, id int, mode string, reassignTo *int) error
	GetDepartmentTree(ctx context.Context, rootID int, maxDepth int) (*models.Department, error)

	CreateEmployee(ctx context.Context, emp *models.Employee) error
	GetEmployeesByDepartmentID(ctx context.Context, deptID int) ([]models.Employee, error)

	IsNameUniqueWithinParent(ctx context.Context, name string, parentId *int, excludeID int) (bool, error)
	IsDescendant(ctx context.Context, parentID, childID int) (bool, error)
}
