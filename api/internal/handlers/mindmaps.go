package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mymindmap/api/auth"
	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"
)

type MindMapHandler struct {
	mindMapRepo *repository.MindMapRepository
	authService *auth.AuthService
	logger      *log.Logger
}

func NewMindMapHandler(mindMapRepo *repository.MindMapRepository, authService *auth.AuthService, logger *log.Logger) *MindMapHandler {
	return &MindMapHandler{
		mindMapRepo: mindMapRepo,
		authService: authService,
		logger:      logger,
	}
}

func (h *MindMapHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/mindmaps", h.authService.AuthMiddleware(h.handleMindMaps))       // GET list, POST create
	mux.HandleFunc("/api/mindmaps/", h.authService.AuthMiddleware(h.handleSingleMindMap)) // GET, PUT, DELETE by id
}

// --- Handlers ---

// handleMindMaps -> /api/mindmaps
func (h *MindMapHandler) handleMindMaps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetMindMaps(w, r)
	case http.MethodPost:
		h.CreateMindMap(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSingleMindMap -> /api/mindmaps/{id}
func (h *MindMapHandler) handleSingleMindMap(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/mindmaps/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid mindmap id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetMindMap(w, r, id)
	case http.MethodPut:
		h.UpdateMindMap(w, r, id)
	case http.MethodDelete:
		h.DeleteMindMap(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetMindMaps - список
func (h *MindMapHandler) GetMindMaps(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	maps, err := h.mindMapRepo.GetMindMapsByUser(r.Context(), user.UserID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, maps)
}

// GetMindMap - один mindmap
func (h *MindMapHandler) GetMindMap(w http.ResponseWriter, r *http.Request, id int) {
	mindmap, err := h.mindMapRepo.GetMindMapByID(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if mindmap == nil {
		h.respondError(w, http.StatusNotFound, "mindmap not found")
		return
	}

	user := auth.GetUserFromContext(r.Context())
	if user == nil || mindmap.UserID != user.UserID && user.Role != "admin" {
		h.respondError(w, http.StatusForbidden, "forbidden")
		return
	}

	h.respondJSON(w, http.StatusOK, mindmap)
}

// CreateMindMap
func (h *MindMapHandler) CreateMindMap(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
		Data  string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Title == "" {
		h.respondError(w, http.StatusBadRequest, "title required")
		return
	}

	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	mindmap := &models.MindMap{
		Title:  req.Title,
		Data:   req.Data,
		UserID: user.UserID,
	}
	if err := h.mindMapRepo.CreateMindMap(r.Context(), mindmap); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, mindmap)
}

// UpdateMindMap
func (h *MindMapHandler) UpdateMindMap(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		Title string `json:"title"`
		Data  string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	mindmap, err := h.mindMapRepo.GetMindMapByID(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if mindmap == nil {
		h.respondError(w, http.StatusNotFound, "mindmap not found")
		return
	}
	if mindmap.UserID != user.UserID && user.Role != "admin" {
		h.respondError(w, http.StatusForbidden, "forbidden")
		return
	}

	mindmap.Title, mindmap.Data = req.Title, req.Data
	if err := h.mindMapRepo.UpdateMindMap(r.Context(), mindmap); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, mindmap)
}

// DeleteMindMap
func (h *MindMapHandler) DeleteMindMap(w http.ResponseWriter, r *http.Request, id int) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	mindmap, err := h.mindMapRepo.GetMindMapByID(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if mindmap == nil {
		h.respondError(w, http.StatusNotFound, "mindmap not found")
		return
	}
	if mindmap.UserID != user.UserID && user.Role != "admin" {
		h.respondError(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.mindMapRepo.DeleteMindMap(r.Context(), id); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

// --- Helpers ---

func (h *MindMapHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("json encode error: %v", err)
	}
}

func (h *MindMapHandler) respondError(w http.ResponseWriter, status int, msg string) {
	h.respondJSON(w, status, map[string]string{"error": msg})
}
