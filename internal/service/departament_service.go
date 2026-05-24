package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/DEELAGRA/org-struct-api/internal/repository"
	"github.com/DEELAGRA/org-struct-api/models"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("conflict")         // например, цикл в дереве
	ErrValidation     = errors.New("validation error") // невалидные данные
	ErrParentNotExist = errors.New("parent department does not exist")
)

type DepartmentService struct {
	repo repository.DepartmentRepository
}

func (s *DepartmentService) CreateDepartment(ctx context.Context, name string, parentID *int) (*models.Department, error) {
	name = strings.TrimSpace(name)
	if len(name) < 1 || len(name) > 200 {
		return nil, ErrValidation
	}

	if parentID != nil {
		parent, err := s.repo.GetDepartment(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, ErrParentNotExist
		}
	}

	unique, err := s.repo.IsNameUniqueWithinParent(ctx, name, parentID, 0)
	if err != nil {
		return nil, err
	}
	if !unique {
		return nil, ErrConflict
	}

	dept := &models.Department{
		Name:     name,
		ParentID: parentID,
	}
	if err := s.repo.CreateDepartment(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *DepartmentService) CreateEmployee(ctx context.Context, deptId int, fullname, position string, hiredAt *string) (*models.Employee, error) {
	fullname = strings.TrimSpace(fullname)
	position = strings.TrimSpace(position)
	if len(fullname) < 1 || len(fullname) > 200 || len(position) < 1 || len(position) > 200 {
		return nil, ErrValidation
	}

	dept, err := s.repo.GetDepartment(ctx, deptId)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, ErrNotFound
	}

	var hiredAtPtr *time.Time
	if hiredAt != nil && *hiredAt != "" {
		t, err := time.Parse("2006-01-02", *hiredAt)
		if err != nil {
			return nil, ErrValidation
		}
		hiredAtPtr = &t
	}

	emp := &models.Employee{
		DepartmentID: deptId,
		FullName:     fullname,
		Position:     position,
		HiredAt:      hiredAtPtr,
	}
	if err := s.repo.CreateEmployee(ctx, emp); err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *DepartmentService) MoveDepartment(ctx context.Context, id int, name *string, parentID *int) (*models.Department, error) {
	dept, err := s.repo.GetDepartment(ctx, id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, ErrNotFound
	}

	if name != nil {
		trimmed := strings.TrimSpace(*name)
		if len(trimmed) < 1 || len(trimmed) > 200 {
			return nil, ErrValidation
		}
		dept.Name = trimmed
	}

	if parentID != nil {
		if *parentID == id {
			return nil, ErrConflict
		}

		isDesc, err := s.repo.IsDescendant(ctx, id, *parentID)
		if err != nil {
			return nil, err
		}
		if isDesc {
			return nil, ErrParentNotExist
		}
		dept.ParentID = parentID

	}

	if name != nil || parentID != nil {
		effectiveName := dept.Name
		unique, err := s.repo.IsNameUniqueWithinParent(ctx, effectiveName, dept.ParentID, id)
		if err != nil {
			return nil, err
		}
		if !unique {
			return nil, ErrConflict
		}

	}

	return dept, nil
}

func (s *DepartmentService) DeleteDepartment(ctx context.Context, id int, mode string, reassignTo *int) error {
	dept, err := s.repo.GetDepartment(ctx, id)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrNotFound
	}
	switch mode {
	case "cascade":
		return s.repo.DeleteDepartment(ctx, id, mode, nil)

	case "reassign":
		if reassignTo == nil {
			return ErrValidation
		}
		targetDept, err := s.repo.GetDepartment(ctx, *reassignTo)
		if err != nil {
			return err
		}
		if targetDept == nil {
			return ErrNotFound
		}
		return s.repo.DeleteDepartment(ctx, id, mode, reassignTo)
	default:
		return ErrValidation
	}
}
