package database

import (
	"cmp"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

// spec is the on-disk shape of a database template.
type spec struct {
	Backup          *backupSpec `json:"backup,omitzero"`
	Restore         *execSpec   `json:"restore,omitzero"`
	Name            string      `json:"name"`
	Repository      string      `json:"repository"`
	Tag             string      `json:"tag"`
	DataDir         string      `json:"data_dir"`
	URL             string      `json:"url"`
	Env             []string    `json:"env,omitzero"`
	Command         []string    `json:"command,omitzero"`
	Client          []string    `json:"client,omitzero"`
	Healthcheck     []string    `json:"healthcheck,omitzero"`
	UnsupportedArch []string    `json:"unsupported_arch,omitzero"`
	Port            int         `json:"port"`
}

type backupSpec struct {
	Extension string   `json:"extension"`
	Command   []string `json:"command"`
	Env       []string `json:"env,omitzero"`
}

type execSpec struct {
	Command []string `json:"command"`
	Env     []string `json:"env,omitzero"`
}

// Definition is a compiled database template.
type Definition struct {
	backup          *backupTemplate
	restore         *execTemplate
	url             *template.Template
	name            string
	repository      string
	tag             string
	dataDir         string
	env             []*template.Template
	command         []*template.Template
	client          []*template.Template
	healthcheck     []*template.Template
	unsupportedArch []string
	port            int
}

// values are the fields templates can reference.
type values struct {
	Env      map[string]string
	Password string
	Host     string
}

// Name returns the database name, such as "postgres".
func (d *Definition) Name() string {
	return d.name
}

// Image returns the image reference for the tag, or the default image when the tag is empty.
func (d *Definition) Image(tag string) string {
	return d.repository + ":" + cmp.Or(tag, d.tag)
}

// Port returns the port the database listens on inside its container.
func (d *Definition) Port() int {
	return d.port
}

// DataDir returns the directory inside the container where the database stores its data.
func (d *Definition) DataDir() string {
	return d.dataDir
}

// Supports reports whether the database image runs on the given GOARCH.
func (d *Definition) Supports(arch string) bool {
	return !slices.Contains(d.unsupportedArch, arch)
}

// Env renders the container environment.
func (d *Definition) Env(password string) ([]string, error) {
	return renderAll(d.env, &values{Password: password})
}

// Command renders the container command, or returns nil to use the image default.
func (d *Definition) Command(password string) ([]string, error) {
	return renderAll(d.command, &values{Password: password})
}

// Healthcheck renders the readiness probe, or returns nil when the template has none.
func (d *Definition) Healthcheck(password string) ([]string, error) {
	return renderAll(d.healthcheck, &values{Password: password})
}

// Client renders the interactive client command, or returns nil when the template has none.
func (d *Definition) Client(password string, env map[string]string) ([]string, error) {
	return renderAll(d.client, &values{Password: password, Env: env})
}

// URL renders the address clients use to connect.
func (d *Definition) URL(host string, port int, password string) (string, error) {
	return render(d.url, &values{Password: password, Host: net.JoinHostPort(host, strconv.Itoa(port))})
}

func (s *spec) compile() (*Definition, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	def := &Definition{
		name:            s.Name,
		repository:      s.Repository,
		tag:             s.Tag,
		dataDir:         s.DataDir,
		unsupportedArch: s.UnsupportedArch,
		port:            s.Port,
	}

	if err := def.parseTemplates(s); err != nil {
		return nil, err
	}

	return def, nil
}

func (d *Definition) parseTemplates(s *spec) error {
	var err error

	if d.env, err = parseAll(s.Name+".env", s.Env); err != nil {
		return err
	}

	if d.command, err = parseAll(s.Name+".command", s.Command); err != nil {
		return err
	}

	if d.client, err = parseAll(s.Name+".client", s.Client); err != nil {
		return err
	}

	if d.healthcheck, err = parseAll(s.Name+".healthcheck", s.Healthcheck); err != nil {
		return err
	}

	if d.url, err = parse(s.Name+".url", s.URL); err != nil {
		return err
	}

	if d.backup, err = s.Backup.compile(s.Name + ".backup"); err != nil {
		return err
	}

	d.restore, err = s.Restore.compile(s.Name + ".restore")

	return err
}

func (s *spec) validate() error {
	switch {
	case s.Name == "":
		return fmt.Errorf("%w: name is required", ErrInvalidTemplate)
	case s.Repository == "" || s.Tag == "":
		return fmt.Errorf("%w: %s: repository and tag are required", ErrInvalidTemplate, s.Name)
	case s.Port <= 0:
		return fmt.Errorf("%w: %s: port must be positive", ErrInvalidTemplate, s.Name)
	case s.DataDir == "":
		return fmt.Errorf("%w: %s: data_dir is required", ErrInvalidTemplate, s.Name)
	case s.URL == "":
		return fmt.Errorf("%w: %s: url is required", ErrInvalidTemplate, s.Name)
	case s.Backup != nil && (s.Backup.Extension == "" || len(s.Backup.Command) == 0):
		return fmt.Errorf("%w: %s: backup needs a command and an extension", ErrInvalidTemplate, s.Name)
	case s.Restore != nil && len(s.Restore.Command) == 0:
		return fmt.Errorf("%w: %s: restore needs a command", ErrInvalidTemplate, s.Name)
	}

	return nil
}

func parse(name, text string) (*template.Template, error) {
	t, err := template.New(name).Option("missingkey=error").Parse(text)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidTemplate, err)
	}

	return t, nil
}

func parseAll(name string, texts []string) ([]*template.Template, error) {
	out := make([]*template.Template, 0, len(texts))

	for i, text := range texts {
		t, err := parse(name+"["+strconv.Itoa(i)+"]", text)
		if err != nil {
			return nil, err
		}

		out = append(out, t)
	}

	return out, nil
}

func render(t *template.Template, v *values) (string, error) {
	var b strings.Builder

	if err := t.Execute(&b, v); err != nil {
		return "", fmt.Errorf("render %s: %w", t.Name(), err)
	}

	return b.String(), nil
}

// a nil result keeps docker on the image default
func renderAll(ts []*template.Template, v *values) ([]string, error) {
	if len(ts) == 0 {
		return nil, nil
	}

	out := make([]string, 0, len(ts))

	for _, t := range ts {
		s, err := render(t, v)
		if err != nil {
			return nil, err
		}

		out = append(out, s)
	}

	return out, nil
}
