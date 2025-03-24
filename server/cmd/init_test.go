package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
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

// Test for InitDB
func TestInitDB(t *testing.T) {
	// Set dummy values for SUPABASE_URL and SUPABASE_ANON_KEY.
	os.Setenv("SUPABASE_URL", "https://dummy.supabase.co")
	os.Setenv("SUPABASE_ANON_KEY", "dummykey")
	defer os.Unsetenv("SUPABASE_URL")
	defer os.Unsetenv("SUPABASE_ANON_KEY")

	app := &App{}

	// Call InitDB which initializes app.DB.
	InitDB(app)

	if app.DB == nil {
		t.Error("Expected app.DB to be non-nil after InitDB")
	}
}

// Test for loggerMiddleware
func TestLoggerMiddleware(t *testing.T) {
	// Create a bytes.Buffer to capture log output.
	var buf bytes.Buffer
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal("Failed to create Zap logger:", err)
	}
	defer zapLogger.Sync()

	testLogger := zapLogger.Sugar()

	// Create a dummy next handler that writes a response and sets headers.
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers that the middleware expects to log.
		w.Header().Set("Status", "200 OK")
		w.Header().Set("StatusCode", "200")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	})

	// Wrap the dummy handler with the logger middleware.
	middlewareHandler := loggerMiddleware(testLogger)(dummyHandler)

	// Create a test request with a body, query parameters, and headers.
	reqBody := "test payload"
	req := httptest.NewRequest("POST", "http://example.com/test?foo=bar", strings.NewReader(reqBody))
	req.Header.Set("Test-Header", "value")
	// Include an Authorization header to verify it is excluded from logging.
	req.Header.Set("Authorization", "Bearer secret")

	rr := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rr, req)

	// Verify response from the dummy handler.
	if rr.Code != http.StatusOK {
		t.Errorf("Expected response code 200, got %d", rr.Code)
	}
	if body := rr.Body.String(); body != "Hello, World!" {
		t.Errorf("Unexpected response body: %q", body)
	}

	// Now, inspect the log output.
	logOutput := buf.String()
	if !strings.Contains(logOutput, "HTTP Request") {
		t.Errorf("Expected log output to contain 'HTTP Request', got %q", logOutput)
	}
	if !strings.Contains(logOutput, reqBody) {
		t.Errorf("Expected log output to contain request payload %q, got %q", reqBody, logOutput)
	}
	// Verify that the Authorization header is excluded.
	if strings.Contains(logOutput, "Authorization") {
		t.Errorf("Expected log output to not contain 'Authorization' header, got %q", logOutput)
	}
	// Check for query parameter logging.
	if !strings.Contains(logOutput, `"foo":"bar"`) {
		t.Errorf("Expected log output to contain query parameter foo=bar, got %q", logOutput)
	}
}
