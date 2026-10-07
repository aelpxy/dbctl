package hostcpu_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/aelpxy/dbctl/internal/hostcpu"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cpuinfo string
		want    []string
	}{
		{
			name:    "raspberry pi 4",
			cpuinfo: "processor\t: 0\nBogoMIPS\t: 108.00\nFeatures\t: fp asimd evtstrm crc32 cpuid\nCPU implementer\t: 0x41\n",
			want:    []string{"fp", "asimd", "evtstrm", "crc32", "cpuid"},
		},
		{
			name:    "x86",
			cpuinfo: "processor\t: 0\nflags\t\t: fpu vme sse2 avx2\n",
			want:    []string{"fpu", "vme", "sse2", "avx2"},
		},
		{
			name:    "no features",
			cpuinfo: "processor\t: 0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := hostcpu.Parse(strings.NewReader(tt.cpuinfo))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
