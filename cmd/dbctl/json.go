package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/aelpxy/dbctl/internal/docker"
)

var (
	errNoPassword = errors.New("password unknown, the database was created before dbctl stored it")
	errNoPort     = errors.New("database port is not published")
)

type databaseJSON struct {
	Created time.Time  `json:"created,omitzero"`
	Stats   *statsJSON `json:"stats,omitzero"`
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Type    string     `json:"type"`
	Image   string     `json:"image"`
	State   string     `json:"state"`
	Status  string     `json:"status"`
	Health  string     `json:"health,omitzero"`
	Ports   []portJSON `json:"ports,omitzero"`
	Volumes []string   `json:"volumes,omitzero"`
}

type portJSON struct {
	HostIP        string `json:"host_ip"`
	HostPort      string `json:"host_port"`
	ContainerPort string `json:"container_port"`
}

type statsJSON struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsageMB float64 `json:"memory_usage_mb"`
	MemoryLimitMB float64 `json:"memory_limit_mb"`
}

type connectionJSON struct {
	ID       string `json:"id,omitzero"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Image    string `json:"image,omitzero"`
	Host     string `json:"host"`
	Password string `json:"password"`
	URL      string `json:"url"`
	Port     int    `json:"port"`
}

func newDatabaseJSON(db *docker.Database) databaseJSON {
	ports := make([]portJSON, 0, len(db.Ports))

	for _, p := range db.Ports {
		ports = append(ports, portJSON{HostIP: p.HostIP, HostPort: p.HostPort, ContainerPort: p.ContainerPort})
	}

	volumes := make([]string, 0, len(db.Volumes))

	for _, v := range db.Volumes {
		volumes = append(volumes, v.Name)
	}

	return databaseJSON{
		Created: db.Created,
		ID:      db.ID,
		Name:    db.Name,
		Type:    db.Type,
		Image:   db.Image,
		State:   db.State,
		Status:  db.Status,
		Health:  db.Health,
		Ports:   ports,
		Volumes: volumes,
	}
}

func newStatsJSON(s docker.Stats) *statsJSON {
	return &statsJSON{
		CPUPercent:    s.CPUPercent,
		MemoryUsageMB: s.MemoryUsageMB(),
		MemoryLimitMB: s.MemoryLimitMB(),
	}
}

// connectionFor rebuilds the connection details of an existing database from its template.
func connectionFor(def *database.Definition, db *docker.Database) (*connectionJSON, error) {
	if db.Password == "" {
		return nil, errNoPassword
	}

	containerPort := strconv.Itoa(def.Port())

	for _, p := range db.Ports {
		if p.ContainerPort != containerPort {
			continue
		}

		port, err := strconv.Atoi(p.HostPort)
		if err != nil {
			return nil, fmt.Errorf("parse host port %q: %w", p.HostPort, err)
		}

		return newConnection(def, db.Name, p.HostIP, port, db.Password)
	}

	return nil, errNoPort
}

func newConnection(def *database.Definition, name, host string, port int, password string) (*connectionJSON, error) {
	url, err := def.URL(host, port, password)
	if err != nil {
		return nil, fmt.Errorf("render connection string: %w", err)
	}

	return &connectionJSON{
		Name:     name,
		Type:     def.Name(),
		Host:     host,
		Password: password,
		URL:      url,
		Port:     port,
	}, nil
}
