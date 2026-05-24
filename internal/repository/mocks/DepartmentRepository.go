package mocks

import (
	"context"

	"github.com/DEELAGRA/org-struct-api/models"
	"github.com/stretchr/testify/mock"
)

type DepartmentRepository struct {
	mock.Mock
}

func (m *DepartmentRepository) CreateDepartment(ctx context.Context, dept *models.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *DepartmentRepository) GetDepartment(ctx context.Context, id int) (*models.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *DepartmentRepository) UpdateDepartment(ctx context.Context, dept *models.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *DepartmentRepository) DeleteDepartment(ctx context.Context, id int, mode string, reassignTo *int) error {
	args := m.Called(ctx, id, mode, reassignTo)
	return args.Error(0)
}

func (m *DepartmentRepository) CreateEmployee(ctx context.Context, emp *models.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *DepartmentRepository) GetDepartmentTree(ctx context.Context, rootID int, maxDepth int) (*models.Department, error) {
	args := m.Called(ctx, rootID, maxDepth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *DepartmentRepository) GetEmployeesByDepartmentID(ctx context.Context, deptID int) ([]models.Employee, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *DepartmentRepository) IsNameUniqueWithinParent(ctx context.Context, name string, parentID *int, excludeID int) (bool, error) {
	args := m.Called(ctx, name, parentID, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *DepartmentRepository) IsDescendant(ctx context.Context, parentID, childID int) (bool, error) {
	args := m.Called(ctx, parentID, childID)
	return args.Bool(0), args.Error(1)
}
