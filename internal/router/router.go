package router

import (
	"net/http"

	"github.com/DEELAGRA/org-struct-api/internal/handler"
	"github.com/DEELAGRA/org-struct-api/internal/service"
)

func SetupRouter(svc *service.DepartmentService) http.Handler {
	mux := http.NewServeMux()
	deptHandler := handler.NewDepartmentHandler(svc)
	deptHandler.RegisterRoutes(mux)
	return mux
}
