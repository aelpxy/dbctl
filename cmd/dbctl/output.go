package main

import (
	"cmp"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

const (
	spinnerCharset  = 14
	spinnerInterval = 80 * time.Millisecond
	shortIDLength   = 12
	cellPadding     = 3
	stepPrecision   = 100 * time.Millisecond
	panelPaddingY   = 1
	panelPaddingX   = 2
)

const (
	colorAccent  = "#A78BFA"
	colorSuccess = "#4ADE80"
	colorWarning = "#FACC15"
	colorError   = "#F87171"
	colorMuted   = "#8B8B99"
)

const (
	symbolSuccess = "✓"
	symbolError   = "✗"
	symbolStatus  = "●"
	symbolHint    = "›"
)

const (
	formatTable outputFormat = "table"
	formatJSON  outputFormat = "json"
)

var errInvalidOutput = errors.New(`invalid output format, use "table" or "json"`)

// outputFormat is a pflag.Value so invalid formats are rejected while parsing flags.
type outputFormat string

func (f *outputFormat) String() string {
	return string(*f)
}

func (f *outputFormat) Set(v string) error {
	switch outputFormat(v) {
	case formatTable, formatJSON:
		*f = outputFormat(v)

		return nil
	default:
		return fmt.Errorf("%w: %s", errInvalidOutput, v)
	}
}

func (f *outputFormat) Type() string {
	return "format"
}

func (f *outputFormat) isJSON() bool {
	return *f == formatJSON
}

func addOutputFlag(cmd *cobra.Command) *outputFormat {
	f := formatTable
	cmd.Flags().VarP(&f, "output", "o", `Output format: "table" or "json".`)

	return &f
}

func writeJSON[T any](v T) error {
	if err := json.MarshalWrite(os.Stdout, v, jsontext.WithIndent("  ")); err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	fmt.Println()

	return nil
}

func fg(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func accent(s string) string {
	return fg(colorAccent).Bold(true).Render(s)
}

func muted(s string) string {
	return fg(colorMuted).Render(s)
}

func bold(s string) string {
	return lipgloss.NewStyle().Bold(true).Render(s)
}

// printOut writes results to stdout; colors are stripped when piped or NO_COLOR is set.
func printOut(s string) {
	_, _ = lipgloss.Println(s)
}

// printErr writes progress and diagnostics to stderr so stdout stays pipeable.
func printErr(s string) {
	_, _ = lipgloss.Fprintln(os.Stderr, s)
}

func successLine(msg string) string {
	return fg(colorSuccess).Render(symbolSuccess) + " " + msg
}

func errorLine(msg string) string {
	return fg(colorError).Render(symbolError) + " " + msg
}

func hintLine(msg string) string {
	return muted("  " + symbolHint + " " + msg)
}

// step shows a spinner while fn runs, then leaves a ✓ or ✗ line with the elapsed time.
func step(msg string, fn func() error) error {
	s := spinner.New(
		spinner.CharSets[spinnerCharset],
		spinnerInterval,
		spinner.WithColor("magenta"),
		spinner.WithSuffix(" "+msg),
		spinner.WithWriter(os.Stderr),
	)
	start := time.Now()

	s.Start()

	err := fn()

	s.Stop()

	elapsed := muted(time.Since(start).Round(stepPrecision).String())

	if err != nil {
		printErr(errorLine(msg + " " + elapsed))

		return err
	}

	printErr(successLine(msg + " " + elapsed))

	return nil
}

func newTable(headers ...string) *table.Table {
	header := fg(colorMuted).Bold(true).PaddingRight(cellPadding)
	cell := lipgloss.NewStyle().PaddingRight(cellPadding)

	return table.New().
		Border(lipgloss.HiddenBorder()).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderColumn(false).
		BorderHeader(false).
		Headers(upper(headers)...).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return header
			}

			return cell
		})
}

func upper(ss []string) []string {
	out := make([]string, 0, len(ss))

	for _, s := range ss {
		out = append(out, strings.ToUpper(s))
	}

	return out
}

// keyValues renders aligned "key  value" lines with dimmed keys.
func keyValues(pairs [][2]string) string {
	width := 0

	for _, p := range pairs {
		width = max(width, lipgloss.Width(p[0]))
	}

	key := fg(colorMuted).Width(width + cellPadding)
	lines := make([]string, 0, len(pairs))

	for _, p := range pairs {
		lines = append(lines, key.Render(p[0])+p[1])
	}

	return strings.Join(lines, "\n")
}

func status(state string) string {
	color := colorMuted

	switch state {
	case "running":
		color = colorSuccess
	case "exited", "dead":
		color = colorError
	case "paused", "restarting", "created":
		color = colorWarning
	}

	return fg(color).Render(symbolStatus) + " " + state
}

func health(h string) string {
	switch h {
	case "healthy":
		return fg(colorSuccess).Render(h)
	case "unhealthy":
		return fg(colorError).Render(h)
	case "":
		return muted("none")
	default:
		return fg(colorWarning).Render(h)
	}
}

func shortID(id string) string {
	return id[:min(len(id), shortIDLength)]
}

func typeName(name string) string {
	return cmp.Or(name, "unknown")
}

func formatPorts(ports []docker.Port) string {
	if len(ports) == 0 {
		return muted("none")
	}

	out := make([]string, 0, len(ports))

	for _, p := range ports {
		out = append(out, p.HostIP+":"+p.HostPort+muted(" → ")+p.ContainerPort)
	}

	return strings.Join(out, ", ")
}

func formatVolumes(volumes []docker.Volume) string {
	if len(volumes) == 0 {
		return muted("none")
	}

	out := make([]string, 0, len(volumes))

	for _, v := range volumes {
		out = append(out, v.Name)
	}

	return strings.Join(out, ", ")
}

// panel frames a block of output in a rounded, accent-colored border.
func panel(body string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorAccent)).
		Padding(panelPaddingY, panelPaddingX).
		Render(body)
}
