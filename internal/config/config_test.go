package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear any environment variables that might interfere
	os.Unsetenv("DEVSMTP_PORT")
	os.Unsetenv("DEVSMTP_HOST")
	os.Unsetenv("DEVSMTP_DB")
	os.Unsetenv("DEVSMTP_AUTH_REQUIRED")
	os.Unsetenv("DEVSMTP_AUTH_USER")
	os.Unsetenv("DEVSMTP_AUTH_PASS")

	cfg, err := Load("", nil)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Check defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host '0.0.0.0', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 587 {
		t.Errorf("expected default port 587, got %d", cfg.Server.Port)
	}
	if cfg.Database.Path != "./devsmtp.db" {
		t.Errorf("expected default db path './devsmtp.db', got %q", cfg.Database.Path)
	}
	if cfg.Auth.Required != false {
		t.Errorf("expected default auth.required false, got %v", cfg.Auth.Required)
	}
	if cfg.Auth.Username != "" {
		t.Errorf("expected default auth.username '', got %q", cfg.Auth.Username)
	}
	if cfg.Auth.Password != "" {
		t.Errorf("expected default auth.password '', got %q", cfg.Auth.Password)
	}
	if cfg.TLS.Cert != "" {
		t.Errorf("expected default tls.cert '', got %q", cfg.TLS.Cert)
	}
	if cfg.TLS.Key != "" {
		t.Errorf("expected default tls.key '', got %q", cfg.TLS.Key)
	}
}

func TestLoadFromEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("DEVSMTP_HOST", "192.168.1.100")
	os.Setenv("DEVSMTP_PORT", "2525")
	os.Setenv("DEVSMTP_DB", "/tmp/test.db")
	defer func() {
		os.Unsetenv("DEVSMTP_HOST")
		os.Unsetenv("DEVSMTP_PORT")
		os.Unsetenv("DEVSMTP_DB")
	}()

	cfg, err := Load("", nil)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Server.Host != "192.168.1.100" {
		t.Errorf("expected host '192.168.1.100', got %q", cfg.Server.Host)
	}
	if cfg.Database.Path != "/tmp/test.db" {
		t.Errorf("expected db path '/tmp/test.db', got %q", cfg.Database.Path)
	}
}

func TestLoadFromConfigFile(t *testing.T) {
	// Create temp config file
	tmpFile, err := os.CreateTemp("", "devsmtp-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
server:
  host: "10.0.0.1"
  port: 1025

database:
  path: "/var/lib/devsmtp/mail.db"

auth:
  required: true
  username: "testuser"
  password: "testpass"

tls:
  cert: "/etc/ssl/cert.pem"
  key: "/etc/ssl/key.pem"
`

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	tmpFile.Close()

	// Clear environment variables
	os.Unsetenv("DEVSMTP_PORT")
	os.Unsetenv("DEVSMTP_HOST")
	os.Unsetenv("DEVSMTP_DB")

	cfg, err := Load(tmpFile.Name(), nil)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Server.Host != "10.0.0.1" {
		t.Errorf("expected host '10.0.0.1', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 1025 {
		t.Errorf("expected port 1025, got %d", cfg.Server.Port)
	}
	if cfg.Database.Path != "/var/lib/devsmtp/mail.db" {
		t.Errorf("expected db path '/var/lib/devsmtp/mail.db', got %q", cfg.Database.Path)
	}
	if cfg.Auth.Required != true {
		t.Errorf("expected auth.required true, got %v", cfg.Auth.Required)
	}
	if cfg.Auth.Username != "testuser" {
		t.Errorf("expected auth.username 'testuser', got %q", cfg.Auth.Username)
	}
	if cfg.Auth.Password != "testpass" {
		t.Errorf("expected auth.password 'testpass', got %q", cfg.Auth.Password)
	}
	if cfg.TLS.Cert != "/etc/ssl/cert.pem" {
		t.Errorf("expected tls.cert '/etc/ssl/cert.pem', got %q", cfg.TLS.Cert)
	}
	if cfg.TLS.Key != "/etc/ssl/key.pem" {
		t.Errorf("expected tls.key '/etc/ssl/key.pem', got %q", cfg.TLS.Key)
	}
}

func TestEnvOverridesConfigFile(t *testing.T) {
	// Create temp config file
	tmpFile, err := os.CreateTemp("", "devsmtp-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
server:
  host: "10.0.0.1"
  port: 1025

database:
  path: "/var/lib/devsmtp/mail.db"
`

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	tmpFile.Close()

	// Set environment variable to override
	os.Setenv("DEVSMTP_HOST", "override.example.com")
	os.Setenv("DEVSMTP_DB", "/override/path.db")
	defer func() {
		os.Unsetenv("DEVSMTP_HOST")
		os.Unsetenv("DEVSMTP_DB")
	}()

	cfg, err := Load(tmpFile.Name(), nil)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Env should override config file
	if cfg.Server.Host != "override.example.com" {
		t.Errorf("expected host 'override.example.com', got %q", cfg.Server.Host)
	}
	if cfg.Database.Path != "/override/path.db" {
		t.Errorf("expected db path '/override/path.db', got %q", cfg.Database.Path)
	}

	// But port should still come from config file
	if cfg.Server.Port != 1025 {
		t.Errorf("expected port 1025, got %d", cfg.Server.Port)
	}
}

func TestConfigStructFields(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 25,
		},
		Database: DatabaseConfig{
			Path: "/tmp/db.sqlite",
		},
		Auth: AuthConfig{
			Required: true,
			Username: "admin",
			Password: "secret",
		},
		TLS: TLSConfig{
			Cert: "/path/to/cert",
			Key:  "/path/to/key",
		},
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Server.Host mismatch")
	}
	if cfg.Server.Port != 25 {
		t.Errorf("Server.Port mismatch")
	}
	if cfg.Database.Path != "/tmp/db.sqlite" {
		t.Errorf("Database.Path mismatch")
	}
	if cfg.Auth.Required != true {
		t.Errorf("Auth.Required mismatch")
	}
	if cfg.Auth.Username != "admin" {
		t.Errorf("Auth.Username mismatch")
	}
	if cfg.Auth.Password != "secret" {
		t.Errorf("Auth.Password mismatch")
	}
	if cfg.TLS.Cert != "/path/to/cert" {
		t.Errorf("TLS.Cert mismatch")
	}
	if cfg.TLS.Key != "/path/to/key" {
		t.Errorf("TLS.Key mismatch")
	}
}

func TestMissingConfigFileUsesDefaults(t *testing.T) {
	// Clear environment
	os.Unsetenv("DEVSMTP_PORT")
	os.Unsetenv("DEVSMTP_HOST")
	os.Unsetenv("DEVSMTP_DB")

	cfg, err := Load("/nonexistent/config.yaml", nil)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Should fall back to defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 587 {
		t.Errorf("expected default port, got %d", cfg.Server.Port)
	}
}
