package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestResolveVersion(t *testing.T) {
	oldVersion := Version
	defer func() { Version = oldVersion }()

	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "injected version with v",
			version: "v1.2.3",
			want:    "1.2.3",
		},
		{
			name:    "injected version without v",
			version: "2.0.0",
			want:    "2.0.0",
		},
		{
			name:    "dev returns dev",
			version: "dev",
			want:    "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Version = tt.version
			got := resolveVersion()
			if got != tt.want {
				t.Errorf("resolveVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionCmdOutput(t *testing.T) {
	oldVersion := Version
	Version = "1.2.3"
	defer func() { Version = oldVersion }()

	buf := new(bytes.Buffer)
	c := &cobra.Command{}
	c.SetOut(buf)
	versionCmd.Run(c, []string{})

	got := buf.String()
	want := "1.2.3\n"
	if got != want {
		t.Errorf("versionCmd output = %q, want %q", got, want)
	}
}

func TestVersionCmdDevOutput(t *testing.T) {
	oldVersion := Version
	Version = "dev"
	defer func() { Version = oldVersion }()

	buf := new(bytes.Buffer)
	c := &cobra.Command{}
	c.SetOut(buf)
	versionCmd.Run(c, []string{})

	got := buf.String()
	want := "dev\n"
	if got != want {
		t.Errorf("versionCmd output = %q, want %q", got, want)
	}
}
