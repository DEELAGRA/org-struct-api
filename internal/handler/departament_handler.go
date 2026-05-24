package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DEELAGRA/org-struct-api/internal/service"
)

type DepartmentHandler struct {
	svc *service.DepartmentService
}

func NewDepartmentHandler(svc *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{svc: svc}
}

func (h *DepartmentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/departments/", h.departmentsRouter)
}

func extractID(path string) (int, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[0] != "departments" {
		return 0, false
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, false
	}
	return id, true
}

func (h *DepartmentHandler) departmentsRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimRight(r.URL.Path, "/")

	if r.Method == http.MethodPost && path == "/departments" {
		h.createDepartment(w, r)
		return
	}

	id, ok := extractID(path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if strings.HasSuffix(path, "/employees") {
		if r.Method == http.MethodPost {
			h.createEmployee(w, r, id)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getDepartment(w, r, id)
	case http.MethodPatch:
		h.updateDepartment(w, r, id)
	case http.MethodDelete:
		h.deleteDepartment(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DepartmentHandler) createDepartment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		ParentID *int   `json:"parent_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dept, err := h.svc.CreateDepartment(r.Context(), req.Name, req.ParentID)
	if err != nil {
		if errors.Is(err, service.ErrValidation) || errors.Is(err, service.ErrParentNotExist) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if errors.Is(err, service.ErrConflict) {
			http.Error(w, "Department name already exists in the same parent", http.StatusConflict)
			return
		}
		log.Printf("Error creating department: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) createEmployee(w http.ResponseWriter, r *http.Request, deptID int) {
	var req struct {
		FullName string  `json:"full_name"`
		Position string  `json:"position"`
		HiredAt  *string `json:"hired_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	emp, err := h.svc.CreateEmployee(r.Context(), deptID, req.FullName, req.Position, req.HiredAt)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Error creating employee: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func (h *DepartmentHandler) getDepartment(w http.ResponseWriter, r *http.Request, id int) {
	depthStr := r.URL.Query().Get("depth")
	depth := 1
	if depthStr != "" {
		var err error
		depth, err = strconv.Atoi(depthStr)
		if err != nil || depth < 1 || depth > 5 {
			depth = 1
		}
	}
	includeEmp := true
	if r.URL.Query().Get("include_employees") == "false" {
		includeEmp = false
	}

	dept, err := h.svc.GetDepartmentTree(r.Context(), id, depth, includeEmp)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting department: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) updateDepartment(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		Name     *string `json:"name"`
		ParentID *int    `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dept, err := h.svc.MoveDepartment(r.Context(), id, req.Name, req.ParentID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrConflict) {
			http.Error(w, "Move would create a cycle or duplicate name", http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Error updating department: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) deleteDepartment(w http.ResponseWriter, r *http.Request, id int) {
	mode := r.URL.Query().Get("mode")
	if mode != "cascade" && mode != "reassign" {
		http.Error(w, "Invalid mode. Must be 'cascade' or 'reassign'", http.StatusBadRequest)
		return
	}
	var reassignTo *int
	if mode == "reassign" {
		reassignStr := r.URL.Query().Get("reassign_to_department_id")
		if reassignStr == "" {
			http.Error(w, "reassign_to_department_id is required for reassign mode", http.StatusBadRequest)
			return
		}
		val, err := strconv.Atoi(reassignStr)
		if err != nil {
			http.Error(w, "Invalid reassign_to_department_id", http.StatusBadRequest)
			return
		}
		reassignTo = &val
	}

	err := h.svc.DeleteDepartment(r.Context(), id, mode, reassignTo)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Error deleting department: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
