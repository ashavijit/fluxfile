package remote

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		connString string
		expectErr  bool
		user       string
		host       string
	}{
		{
			name:       "valid connection string",
			connString: "deploy@prod.example.com",
			expectErr:  false,
			user:       "deploy",
			host:       "prod.example.com",
		},
		{
			name:       "valid localhost connection",
			connString: "root@localhost",
			expectErr:  false,
			user:       "root",
			host:       "localhost",
		},
		{
			name:       "invalid connection string - no @",
			connString: "invalid",
			expectErr:  true,
		},
		{
			name:       "invalid connection string - empty",
			connString: "",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := New(tt.connString)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if r == nil {
				t.Fatal("Expected Remote instance, got nil")
			}

			if r.user != tt.user {
				t.Errorf("Expected user %s, got %s", tt.user, r.user)
			}

			if r.host != tt.host {
				t.Errorf("Expected host %s, got %s", tt.host, r.host)
			}

			if r.logger == nil {
				t.Error("Expected logger to be initialized")
			}
		})
	}
}

func TestCopyFileNotImplemented(t *testing.T) {
	r, err := New("user@host")
	if err != nil {
		t.Fatalf("Failed to create Remote: %v", err)
	}

	err = r.CopyFile("/local/path", "/remote/path")
	if err == nil {
		t.Error("Expected error for unimplemented CopyFile")
	}

	if err.Error() != "file copy not implemented yet" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestGetSSHAuth(t *testing.T) {
	r, _ := New("user@host")
	auth := r.getSSHAuth()
	_ = auth
}

func TestConnectionStringParsing(t *testing.T) {
	r, err := New("admin@192.168.1.100")
	if err != nil {
		t.Fatalf("Failed to parse IP connection string: %v", err)
	}

	if r.host != "192.168.1.100" {
		t.Errorf("Expected host 192.168.1.100, got %s", r.host)
	}
}

func TestConnectionStringWithPort(t *testing.T) {
	// Note: Current implementation doesn't support port in connection string
	// This test documents current behavior
	r, err := New("user@host:2222")
	if err != nil {
		t.Fatalf("Failed to parse connection string: %v", err)
	}

	// Currently, port is included in host (which may cause issues)
	// This test documents the behavior for future improvement
	if r.host != "host:2222" {
		t.Errorf("Expected host:2222 (current behavior), got %s", r.host)
	}
}
