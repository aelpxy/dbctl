// Package database loads the embedded database templates and renders their container configuration.
package database

import (
	"embed"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed templates/*.json
var templates embed.FS

var (
	ErrUnsupported     = errors.New("unsupported database type")
	ErrInvalidTemplate = errors.New("invalid database template")
)

// Registry holds every known database definition.
type Registry struct {
	defs []*Definition
}

// Load compiles the embedded database templates.
func Load() (*Registry, error) {
	paths, err := fs.Glob(templates, "templates/*.json")
	if err != nil {
		return nil, fmt.Errorf("find templates: %w", err)
	}

	r := &Registry{defs: make([]*Definition, 0, len(paths))}

	for _, path := range paths {
		def, err := loadFile(path)
		if err != nil {
			return nil, err
		}

		if _, err := r.Get(def.Name()); err == nil {
			return nil, fmt.Errorf("%w: %s: duplicate name %s", ErrInvalidTemplate, path, def.Name())
		}

		r.defs = append(r.defs, def)
	}

	return r, nil
}

// Names returns the names of every known database.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.defs))

	for _, def := range r.defs {
		names = append(names, def.Name())
	}

	return names
}

// Get returns the database with the given name.
func (r *Registry) Get(name string) (*Definition, error) {
	for _, def := range r.defs {
		if def.Name() == name {
			return def, nil
		}
	}

	return nil, fmt.Errorf("%w: %s (supported: %s)", ErrUnsupported, name, strings.Join(r.Names(), ", "))
}

// Parse splits a "<type>[:tag]" argument into its database and optional image tag.
func (r *Registry) Parse(arg string) (*Definition, string, error) {
	name, tag, _ := strings.Cut(arg, ":")

	def, err := r.Get(name)
	if err != nil {
		return nil, "", err
	}

	return def, tag, nil
}

// FromImage returns the database whose repository matches an image reference such as "postgres:18-alpine".
func (r *Registry) FromImage(image string) (*Definition, error) {
	repository := repositoryOf(image)

	for _, def := range r.defs {
		if def.repository == repository {
			return def, nil
		}
	}

	return nil, fmt.Errorf("%w: image %s", ErrUnsupported, image)
}

// Resolve returns the database by name, falling back to its image for containers created without a name label.
func (r *Registry) Resolve(name, image string) (*Definition, error) {
	if def, err := r.Get(name); err == nil {
		return def, nil
	}

	return r.FromImage(image)
}

func loadFile(path string) (*Definition, error) {
	f, err := templates.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	defer func() { _ = f.Close() }()

	var s spec

	if err := json.UnmarshalRead(f, &s, json.RejectUnknownMembers(true)); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	def, err := s.compile()
	if err != nil {
		return nil, fmt.Errorf("compile %s: %w", path, err)
	}

	return def, nil
}

// strips the tag while keeping registry ports such as localhost:5000/name
func repositoryOf(image string) string {
	i := strings.LastIndex(image, ":")
	if i < 0 || strings.Contains(image[i:], "/") {
		return image
	}

	return image[:i]
}
