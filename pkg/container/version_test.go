package container

import (
	"testing"
)

func TestDaemonAPIVersionAtLeast(t *testing.T) {
	cases := []struct {
		have     string
		required string
		want     bool
	}{
		{"1.44", "1.44", true},
		{"1.45", "1.44", true},
		{"1.43", "1.44", false},
		{"1.9", "1.10", false},
		{"1.10", "1.9", true},
		{"2.0", "1.44", true},
		{"1.44", "2.0", false},
	}
	for _, tc := range cases {
		got := daemonAPIVersionAtLeast(tc.have, tc.required)
		if got != tc.want {
			t.Errorf("daemonAPIVersionAtLeast(%q, %q) = %v, want %v", tc.have, tc.required, got, tc.want)
		}
	}
}
