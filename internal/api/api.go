// Package api serves a read-only JSON view of dbctl databases over HTTP.
package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"log/slog"
	"net/http"

	"github.com/aelpxy/dbctl/internal/docker"
)

// Store provides the databases served by the API.
type Store interface {
	List(ctx context.Context) ([]docker.Database, error)
	Inspect(ctx context.Context, id string) (*docker.Database, error)
	Stats(ctx context.Context, id string) (docker.Stats, error)
}

type response[T any] struct {
	Message string `json:"message,omitzero"`
	Data    []T    `json:"data,omitzero"`
	Success bool   `json:"success"`
}

type summary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	State   string `json:"state"`
	Image   string `json:"image"`
	Created int64  `json:"created"`
}

type detail struct {
	Ports    map[string]string `json:"ports,omitzero"`
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Image    string            `json:"image"`
	Status   string            `json:"status"`
	Volumes  []volume          `json:"volumes,omitzero"`
	Memory   memory            `json:"memory_usage"`
	CPUUsage float64           `json:"cpu_usage"`
}

type memory struct {
	Used  float64 `json:"used_mb"`
	Total float64 `json:"total_mb"`
}

type volume struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type handler struct {
	store  Store
	logger *slog.Logger
}

// NewHandler returns the API routes backed by store.
func NewHandler(store Store, logger *slog.Logger) *http.ServeMux {
	h := &handler{store: store, logger: logger}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", h.healthcheck)
	mux.HandleFunc("GET /databases", h.list)
	mux.HandleFunc("GET /databases/{id}", h.retrieve)

	return mux
}

func (h *handler) healthcheck(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, h.logger, http.StatusOK, response[struct{}]{Success: true, Message: "dbctl API is ready to serve you!"})
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	dbs, err := h.store.List(r.Context())
	if err != nil {
		h.logger.Error("list databases", "error", err)
		h.fail(w, http.StatusInternalServerError, "Failed to list containers.")

		return
	}

	data := make([]summary, 0, len(dbs))

	for i := range dbs {
		data = append(data, newSummary(&dbs[i]))
	}

	resp := response[summary]{Success: true, Data: data}

	if len(data) == 0 {
		resp.Message = "No databases found."
	}

	writeJSON(w, h.logger, http.StatusOK, resp)
}

func (h *handler) retrieve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	db, err := h.store.Inspect(r.Context(), id)
	if err != nil {
		h.inspectFailed(w, id, err)

		return
	}

	stats, err := h.store.Stats(r.Context(), db.ID)
	if err != nil {
		h.logger.Error("get database stats", "database_id", db.ID, "error", err)
		h.fail(w, http.StatusInternalServerError, "Error getting container stats.")

		return
	}

	writeJSON(w, h.logger, http.StatusOK, response[detail]{Success: true, Data: []detail{newDetail(db, stats)}})
}

func (h *handler) inspectFailed(w http.ResponseWriter, id string, err error) {
	if errors.Is(err, docker.ErrNotFound) || errors.Is(err, docker.ErrNotManaged) {
		h.fail(w, http.StatusNotFound, "Database not found.")

		return
	}

	h.logger.Error("inspect database", "database_id", id, "error", err)
	h.fail(w, http.StatusInternalServerError, "Error inspecting container.")
}

func (h *handler) fail(w http.ResponseWriter, status int, message string) {
	writeJSON(w, h.logger, status, response[struct{}]{Message: message})
}

func writeJSON[T any](w http.ResponseWriter, logger *slog.Logger, status int, body response[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.MarshalWrite(w, body); err != nil {
		logger.Error("write response", "error", err)
	}
}

func newSummary(db *docker.Database) summary {
	return summary{
		ID:      db.ID,
		Name:    db.Name,
		Type:    db.Type,
		Status:  db.Status,
		State:   db.State,
		Image:   db.Image,
		Created: db.Created.Unix(),
	}
}

func newDetail(db *docker.Database, stats docker.Stats) detail {
	ports := make(map[string]string, len(db.Ports))

	for _, p := range db.Ports {
		ports[p.ContainerPort] = p.HostPort
	}

	volumes := make([]volume, 0, len(db.Volumes))

	for _, v := range db.Volumes {
		volumes = append(volumes, volume{Source: v.Source, Target: v.Target})
	}

	return detail{
		Ports:    ports,
		ID:       db.ID,
		Name:     db.Name,
		Type:     db.Type,
		Image:    db.Image,
		Status:   db.State,
		Volumes:  volumes,
		Memory:   memory{Used: stats.MemoryUsageMB(), Total: stats.MemoryLimitMB()},
		CPUUsage: stats.CPUPercent,
	}
}
