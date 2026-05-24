package service

import (
	"context"
	"testing"

	"github.com/DEELAGRA/org-struct-api/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateDepartment_EmptyName(t *testing.T) {

	repo := new(mocks.DepartmentRepository)
	svc := NewDepartmentService(repo)
	ctx := context.Background()

	dept, err := svc.CreateDepartment(ctx, "   ", nil)

	assert.Nil(t, dept)
	assert.ErrorIs(t, err, ErrValidation)
	repo.AssertNotCalled(t, "CreateDepartment")
}
