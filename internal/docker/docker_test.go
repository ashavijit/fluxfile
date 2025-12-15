package docker

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		expected string
	}{
		{
			name:     "custom image",
			image:    "golang:1.22",
			expected: "golang:1.22",
		},
		{
			name:     "empty image defaults to alpine",
			image:    "",
			expected: "alpine:latest",
		},
		{
			name:     "node image",
			image:    "node:20-alpine",
			expected: "node:20-alpine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New(tt.image)
			if d == nil {
				t.Fatal("Expected Docker instance, got nil")
			}
			if d.image != tt.expected {
				t.Errorf("Expected image %s, got %s", tt.expected, d.image)
			}
			if d.logger == nil {
				t.Error("Expected logger to be initialized")
			}
		})
	}
}

func TestIsAvailable(t *testing.T) {
	d := New("alpine:latest")

	// This test checks if Docker is available on the system
	// It's a smoke test - result depends on whether Docker is installed
	available := d.IsAvailable()

	// We just check that the function doesn't panic
	_ = available
}

func TestDockerCommandArgs(t *testing.T) {
	// Test that Docker struct is properly initialized
	d := New("golang:1.22")

	if d.image != "golang:1.22" {
		t.Errorf("Expected image golang:1.22, got %s", d.image)
	}
}
