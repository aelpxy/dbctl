// Package hostcpu reports the instruction set features of the host processor.
package hostcpu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const cpuinfoPath = "/proc/cpuinfo"

// Features returns the CPU feature flags listed in /proc/cpuinfo ("Features" on arm, "flags" on x86).
// ok is false when they cannot be determined, e.g. outside linux, and callers should assume support.
func Features() (features []string, ok bool) {
	f, err := os.Open(cpuinfoPath)
	if err != nil {
		return nil, false
	}

	defer func() { _ = f.Close() }()

	features, err = Parse(f)
	if err != nil || len(features) == 0 {
		return nil, false
	}

	return features, true
}

// Parse returns the feature flags of the first processor listed in cpuinfo-formatted r.
func Parse(r io.Reader) ([]string, error) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		key, value, found := strings.Cut(s.Text(), ":")
		if !found {
			continue
		}

		switch strings.TrimSpace(key) {
		case "Features", "flags":
			return strings.Fields(value), nil
		}
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", cpuinfoPath, err)
	}

	return nil, nil
}
