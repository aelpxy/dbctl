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

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/briandowns/spinner"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	spinnerCharset  = 11
	spinnerInterval = 100 * time.Millisecond
	shortIDLength   = 12
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

// progress goes to stderr so stdout stays clean for json and piping
func startSpinner(message string) *spinner.Spinner {
	s := spinner.New(
		spinner.CharSets[spinnerCharset],
		spinnerInterval,
		spinner.WithColor("green"),
		spinner.WithSuffix(" "+message),
		spinner.WithWriter(os.Stderr),
	)
	s.Start()

	return s
}

func progressf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func newTable(header ...string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Right: true})
	table.SetCenterSeparator("|")

	return table
}

func shortID(id string) string {
	return id[:min(len(id), shortIDLength)]
}

func typeName(name string) string {
	return cmp.Or(name, "unknown")
}

func statusColor(state string) tablewriter.Colors {
	switch state {
	case "running":
		return tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold}
	case "exited", "dead":
		return tablewriter.Colors{tablewriter.FgRedColor, tablewriter.Bold}
	case "paused", "restarting":
		return tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold}
	default:
		return tablewriter.Colors{}
	}
}

func formatPorts(ports []docker.Port) string {
	if len(ports) == 0 {
		return "none"
	}

	out := make([]string, 0, len(ports))

	for _, p := range ports {
		out = append(out, p.HostIP+":"+p.HostPort+"->"+p.ContainerPort)
	}

	return strings.Join(out, ", ")
}

func formatVolumes(volumes []docker.Volume) string {
	if len(volumes) == 0 {
		return "none"
	}

	out := make([]string, 0, len(volumes))

	for _, v := range volumes {
		out = append(out, v.Name)
	}

	return strings.Join(out, ", ")
}
