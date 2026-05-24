package repository

import (
	"context"
	"errors"

	"github.com/DEELAGRA/org-struct-api/models"
	"gorm.io/gorm"
)

type departmentRepo struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepo{db: db}
}

func (r *departmentRepo) CreateDepartment(ctx context.Context, dept *models.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

func (r *departmentRepo) GetDepartment(ctx context.Context, id int) (*models.Department, error) {
	var dept models.Department
	err := r.db.WithContext(ctx).Preload("Parent").First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepo) UpdateDepartment(ctx context.Context, dept *models.Department) error {
	return r.db.WithContext(ctx).Save(dept).Error
}

func (r *departmentRepo) DeleteDepartment(ctx context.Context, id int, mode string, reassignTo *int) error {
	return r.db.WithContext(ctx).Delete(&models.Department{}, id).Error
}

func (r *departmentRepo) CreateEmployee(ctx context.Context, emp *models.Employee) error {
	var dept models.Department
	if err := r.db.WithContext(ctx).First(&dept, emp.DepartmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("department not found")
		}
		return err
	}
	return r.db.WithContext(ctx).Create(emp).Error
}

func (r *departmentRepo) IsNameUniqueWithinParent(ctx context.Context, name string, parentId *int, excludeID int) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Department{}).
		Where("name = ? AND parent_id IS NOT DISTINCT FROM ?", name, parentId)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *departmentRepo) IsDescendant(ctx context.Context, parentId, childID int) (bool, error) {
	currentID := childID
	for {
		var dept models.Department
		err := r.db.WithContext(ctx).Select("parent_id").First(&dept, currentID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, nil
		}
		if dept.ParentID == nil {
			return false, nil
		}
		if *dept.ParentID == parentId {
			return true, nil
		}
		currentID = *dept.ParentID
	}
}

func (r *departmentRepo) GetDepartmentTree(ctx context.Context, rootID int, maxDepth int) (*models.Department, error) {
	var root models.Department
	if err := r.db.WithContext(ctx).First(&root, rootID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if err := r.loadChildren(ctx, &root, 1, maxDepth); err != nil {
		return nil, err
	}
	return &root, nil
}

func (r *departmentRepo) loadChildren(ctx context.Context, dept *models.Department, currentDepth, maxDepth int) error {
	if currentDepth > maxDepth {
		return nil
	}
	var children []models.Department
	if err := r.db.WithContext(ctx).Where("parent_id = ?", dept.ID).Find(&children).Error; err != nil {
		return err
	}
	for i := range children {
		if err := r.loadChildren(ctx, &children[i], currentDepth+1, maxDepth); err != nil {
			return err
		}
	}
	dept.Children = children
	return nil
}

func (r *departmentRepo) GetEmployeesByDepartmentID(ctx context.Context, deptID int) ([]models.Employee, error) {
	var emps []models.Employee
	err := r.db.WithContext(ctx).Where("department_id = ?", deptID).Find(&emps).Error
	return emps, err
}
