package main

import (
	"os"
	"testing"
)

// Test for InitConfig
func TestInitConfig(t *testing.T) {
	// Create a temporary .env file.
	envContent := "TEST_VAR=hello\n"
	const envFile = ".env"
	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to write temporary .env file: %v", err)
	}
	// Ensure the temporary file is removed after test.
	defer os.Remove(envFile)

	// Call InitConfig which uses godotenv.Load().
	InitConfig()

	// Verify that the environment variable from the .env file is loaded.
	if val := os.Getenv("TEST_VAR"); val != "hello" {
		t.Errorf("Expected TEST_VAR to be 'hello', got %q", val)
	}
}

// Test for InitServer
func TestInitServer(t *testing.T) {
	// Unset API_PORT to force the default value.
	os.Unsetenv("API_PORT")
	app := &App{}

	// Call InitServer which registers routes and assigns a server to app.Srv.
	InitServer(app)

	if app.Srv == nil {
		t.Fatal("Expected app.Srv to be non-nil after InitServer")
	}

	// Default port should be 3000 if API_PORT is not set.
	expectedAddr := "0.0.0.0:3000"
	if app.Srv.Addr != expectedAddr {
		t.Errorf("Expected server address to be %q, got %q", expectedAddr, app.Srv.Addr)
	}

	// Ensure the server handler is set.
	if app.Srv.Handler == nil {
		t.Error("Expected server handler to be non-nil")
	}
}
