package database_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/aelpxy/dbctl/internal/database"
)

func load(t *testing.T) *database.Registry {
	t.Helper()

	r, err := database.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	return r
}

func get(t *testing.T, r *database.Registry, name string) *database.Definition {
	t.Helper()

	def, err := r.Get(name)
	if err != nil {
		t.Fatalf("Get(%q) error = %v", name, err)
	}

	return def
}

// containerEnv renders the template env the way docker reports it back on inspect.
func containerEnv(t *testing.T, def *database.Definition, password string) map[string]string {
	t.Helper()

	vars, err := def.Env(password)
	if err != nil {
		t.Fatalf("Env() error = %v", err)
	}

	env := make(map[string]string, len(vars))

	for _, v := range vars {
		key, value, _ := strings.Cut(v, "=")
		env[key] = value
	}

	return env
}

func TestLoad(t *testing.T) {
	t.Parallel()

	r := load(t)

	want := []string{
		"clickhouse", "couchdb", "dragonfly", "keydb", "mariadb", "meilisearch",
		"mongo", "mysql", "pgvector", "postgres", "redis", "valkey",
	}
	if got := r.Names(); !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

func TestRenderAll(t *testing.T) {
	t.Parallel()

	r := load(t)

	for _, name := range r.Names() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			def := get(t, r, name)
			env := containerEnv(t, def, "secret")

			if err := errors.Join(renderErrors(def, env)...); err != nil {
				t.Error(err)
			}
		})
	}
}

// renderErrors renders every template of def; errors name the failing template.
func renderErrors(def *database.Definition, env map[string]string) []error {
	_, commandErr := def.Command("secret")
	_, healthcheckErr := def.Healthcheck("secret")
	_, clientErr := def.Client("secret", env)
	_, urlErr := def.URL("127.0.0.1", 1234, "secret")

	errs := []error{commandErr, healthcheckErr, clientErr, urlErr}

	if def.CanBackup() {
		_, err := def.Backup(env)
		errs = append(errs, err)
	}

	if def.CanRestore() {
		_, err := def.Restore(env)
		errs = append(errs, err)
	}

	return errs
}

func TestParse(t *testing.T) {
	t.Parallel()

	r := load(t)

	tests := []struct {
		wantErr  error
		name     string
		input    string
		wantName string
		wantTag  string
	}{
		{name: "type only", input: "postgres", wantName: "postgres"},
		{name: "type and tag", input: "redis:latest", wantName: "redis", wantTag: "latest"},
		{name: "unsupported", input: "sqlite", wantErr: database.ErrUnsupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			def, tag, err := r.Parse(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Parse(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}

			if err != nil {
				return
			}

			if def.Name() != tt.wantName || tag != tt.wantTag {
				t.Errorf("Parse(%q) = %q, %q, want %q, %q", tt.input, def.Name(), tag, tt.wantName, tt.wantTag)
			}
		})
	}
}

func TestFromImage(t *testing.T) {
	t.Parallel()

	r := load(t)

	tests := []struct {
		image string
		want  string
	}{
		{image: "postgres:17.0-alpine3.19", want: "postgres"},
		{image: "pgvector/pgvector:pg18", want: "pgvector"},
		{image: "getmeili/meilisearch:v1.10.3", want: "meilisearch"},
		{image: "clickhouse/clickhouse-server", want: "clickhouse"},
		{image: "docker.dragonflydb.io/dragonflydb/dragonfly:v1.37.0", want: "dragonfly"},
	}

	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			t.Parallel()

			def, err := r.FromImage(tt.image)
			if err != nil {
				t.Fatalf("FromImage(%q) error = %v", tt.image, err)
			}

			if def.Name() != tt.want {
				t.Errorf("FromImage(%q) = %q, want %q", tt.image, def.Name(), tt.want)
			}
		})
	}

	if _, err := r.FromImage("nginx:latest"); !errors.Is(err, database.ErrUnsupported) {
		t.Errorf("FromImage(nginx) error = %v, want %v", err, database.ErrUnsupported)
	}
}

func TestBackupAndRestore(t *testing.T) {
	t.Parallel()

	r := load(t)
	def := get(t, r, "postgres")
	env := map[string]string{"POSTGRES_USER": "postgres", "POSTGRES_DB": "app"}

	backup, err := def.Backup(env)
	if err != nil {
		t.Fatalf("Backup() error = %v", err)
	}

	if want := []string{"pg_dump", "-U", "postgres", "app"}; !slices.Equal(backup.Command, want) {
		t.Errorf("Backup().Command = %v, want %v", backup.Command, want)
	}

	restore, err := def.Restore(env)
	if err != nil {
		t.Fatalf("Restore() error = %v", err)
	}

	if want := []string{"psql", "-v", "ON_ERROR_STOP=1", "-U", "postgres", "-d", "app"}; !slices.Equal(restore.Command, want) {
		t.Errorf("Restore().Command = %v, want %v", restore.Command, want)
	}

	if _, err := def.Backup(map[string]string{}); err == nil {
		t.Error("Backup() with missing env succeeded, want error")
	}

	redis := get(t, r, "redis")

	if _, err := redis.Backup(nil); !errors.Is(err, database.ErrBackupUnsupported) {
		t.Errorf("Backup() error = %v, want %v", err, database.ErrBackupUnsupported)
	}

	if _, err := redis.Restore(nil); !errors.Is(err, database.ErrRestoreUnsupported) {
		t.Errorf("Restore() error = %v, want %v", err, database.ErrRestoreUnsupported)
	}
}

func TestMissingCPUFeatures(t *testing.T) {
	t.Parallel()

	r := load(t)
	clickhouse := get(t, r, "clickhouse")

	pi4 := strings.Fields("fp asimd evtstrm crc32 cpuid")
	pi5 := strings.Fields("fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm lrcpc dcpop asimddp")

	if got := clickhouse.MissingCPUFeatures("arm64", pi4); len(got) == 0 {
		t.Error("MissingCPUFeatures(arm64, pi4) is empty, want clickhouse rejected on a raspberry pi 4")
	}

	if got := clickhouse.MissingCPUFeatures("arm64", pi5); len(got) != 0 {
		t.Errorf("MissingCPUFeatures(arm64, pi5) = %v, want none", got)
	}

	if got := clickhouse.MissingCPUFeatures("amd64", nil); len(got) != 0 {
		t.Errorf("MissingCPUFeatures(amd64) = %v, want none", got)
	}

	if got := get(t, r, "postgres").MissingCPUFeatures("arm64", pi4); len(got) != 0 {
		t.Errorf("postgres MissingCPUFeatures(arm64, pi4) = %v, want none", got)
	}
}
