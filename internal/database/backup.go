package database

import (
	"errors"
	"fmt"
	"text/template"
)

var (
	ErrBackupUnsupported  = errors.New("backups are not supported for this database")
	ErrRestoreUnsupported = errors.New("restores are not supported for this database")
)

// Exec is a command run inside a database container, streaming the dump over stdout or stdin.
type Exec struct {
	Command []string
	Env     []string
}

type execTemplate struct {
	command []*template.Template
	env     []*template.Template
}

type backupTemplate struct {
	exec      *execTemplate
	extension string
}

// CanBackup reports whether the template defines a backup command.
func (d *Definition) CanBackup() bool {
	return d.backup != nil
}

// CanRestore reports whether the template defines a restore command.
func (d *Definition) CanRestore() bool {
	return d.restore != nil
}

// BackupExtension returns the file extension of the backup format.
func (d *Definition) BackupExtension() string {
	if d.backup == nil {
		return ""
	}

	return d.backup.extension
}

// Backup renders the command that writes a dump to stdout, for a container started with env.
func (d *Definition) Backup(env map[string]string) (Exec, error) {
	if d.backup == nil {
		return Exec{}, fmt.Errorf("%w: %s", ErrBackupUnsupported, d.name)
	}

	return d.backup.exec.render(env)
}

// Restore renders the command that reads a dump from stdin, for a container started with env.
func (d *Definition) Restore(env map[string]string) (Exec, error) {
	if d.restore == nil {
		return Exec{}, fmt.Errorf("%w: %s", ErrRestoreUnsupported, d.name)
	}

	return d.restore.render(env)
}

func (e *execTemplate) render(env map[string]string) (Exec, error) {
	v := &values{Env: env}

	command, err := renderAll(e.command, v)
	if err != nil {
		return Exec{}, err
	}

	vars, err := renderAll(e.env, v)
	if err != nil {
		return Exec{}, err
	}

	return Exec{Command: command, Env: vars}, nil
}

func (b *backupSpec) compile(name string) (*backupTemplate, error) {
	if b == nil {
		return nil, nil
	}

	exec, err := (&execSpec{Command: b.Command, Env: b.Env}).compile(name)
	if err != nil {
		return nil, err
	}

	return &backupTemplate{exec: exec, extension: b.Extension}, nil
}

func (e *execSpec) compile(name string) (*execTemplate, error) {
	if e == nil {
		return nil, nil
	}

	command, err := parseAll(name+".command", e.Command)
	if err != nil {
		return nil, err
	}

	env, err := parseAll(name+".env", e.Env)
	if err != nil {
		return nil, err
	}

	return &execTemplate{command: command, env: env}, nil
}
