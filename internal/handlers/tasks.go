package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

// Handler содержит зависимости HTTP-обработчиков.
type Handler struct{ Store storage.Storage }

// New создаёт новый Handler с переданным хранилищем.
func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// errorResponse — единый формат ошибок для всех эндпоинтов: {"error": "сообщение"}.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON пишет JSON-ответ с нужным статусом и заголовком Content-Type.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// writeError пишет JSON-ответ с ошибкой в едином формате {"error": "..."}.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

// TasksCollection обрабатывает /tasks: GET — список задач, POST — создание задачи.
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	created, err := h.Store.Create(task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// TaskItem обрабатывает /tasks/{id}: GET, PUT, DELETE для конкретной задачи.
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// extractID достаёт числовой ID из пути вида /tasks/{id}.
func extractID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[1] == "" {
		return 0, fmt.Errorf("id is required")
	}
	return strconv.Atoi(parts[1])
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request, id int) {
	task, ok := h.Store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id int) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	updated, err := h.Store.Update(id, task)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.Store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
