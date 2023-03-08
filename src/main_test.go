package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestLockfile() string {
	lockFile := os.TempDir() + "/builderLockFile_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	file, _ := os.Create(lockFile)
	file.Close()
	return lockFile
}

func TestHealthcheck_ShouldSucceed(t *testing.T) {
	r := setupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestBuild_AuthenticateWebhookShouldFailOnMissingBuildCommand(t *testing.T) {
	os.Unsetenv("BUILDER_WEBHOOK_SECRET")

	r := setupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/build", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, "{\"description\":\"Build command missing. Nothing to execute.\"}", w.Body.String())
}

func TestBuild_AuthenticateWebhookShouldFailOnInvalidSecret(t *testing.T) {
	os.Setenv("BUILDER_WEBHOOK_SECRET", "foo")
	defer os.Unsetenv("BUILDER_WEBHOOK_SECRET")

	r := setupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/build", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "{\"description\":\"Authentication failed.\"}", w.Body.String())
}

func TestBuild_AuthenticateWebhookShouldFailOnBuildCommandMissing(t *testing.T) {
	os.Setenv("BUILDER_WEBHOOK_SECRET", "foo")
	defer os.Unsetenv("BUILDER_WEBHOOK_SECRET")

	os.Setenv("BUILDER_WEBHOOK_SECRET_HEADER", "x-build-secret")
	defer os.Unsetenv("BUILDER_WEBHOOK_SECRET_HEADER")

	r := setupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/build", nil)
	req.Header.Add("x-build-secret", "foo")

	r.ServeHTTP(w, req)

	assert.Equal(t, 501, w.Code)
	assert.Equal(t, "{\"description\":\"Build command missing. Nothing to execute.\"}", w.Body.String())
}

func TestBuild_AuthenticateWebhookShouldFailOnLockfileExisting(t *testing.T) {
	lockFilePath := createTestLockfile()
	defer os.Remove(lockFilePath)

	os.Setenv("BUILDER_WEBHOOK_SECRET", "foo")
	defer os.Unsetenv("BUILDER_WEBHOOK_SECRET")

	os.Setenv("BUILDER_WEBHOOK_SECRET_HEADER", "x-build-secret")
	defer os.Unsetenv("BUILDER_WEBHOOK_SECRET_HEADER")

	os.Setenv("BUILDER_EXEC_PATH", "/dev/null")
	defer os.Unsetenv("BUILDER_EXEC_PATH")

	r := setupRoutes()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/build", nil)
	req.Header.Add("x-build-secret", "foo")

	r.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, "{\"description\":\"Lockfile found. Build in progress. Please wait a moment.\"}", w.Body.String())
}

func TestBuild_AuthenticateWebhookShouldSucceedWithoutSecret(t *testing.T) {
	// Create a random lockfile path for this test
	os.Setenv("BUILDER_LOCKFILE_PATH", os.TempDir()+"/builderLockFile_"+strconv.FormatInt(time.Now().UnixNano(), 10)+"_1")
	defer os.Unsetenv("BUILDER_LOCKFILE_PATH")

	os.Unsetenv("BUILDER_WEBHOOK_SECRET")

	os.Setenv("BUILDER_EXEC_PATH", "echo")
	defer os.Unsetenv("BUILDER_EXEC_PATH")

	r := setupRoutes()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/build", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestBuild_AuthenticateWebhookShouldSucceedWithSecret(t *testing.T) {
	// Create a random lockfile path for this test
	os.Setenv("BUILDER_LOCKFILE_PATH", os.TempDir()+"/builderLockFile_"+strconv.FormatInt(time.Now().UnixNano(), 10)+"_2")
	defer os.Unsetenv("BUILDER_LOCKFILE_PATH")

	os.Setenv("BUILDER_WEBHOOK_SECRET", "foo")
	defer os.Unsetenv("BUILDER_WEBHOOK_SECRET")

	os.Setenv("BUILDER_EXEC_PATH", "echo")
	defer os.Unsetenv("BUILDER_EXEC_PATH")

	os.Setenv("BUILDER_WEBHOOK_SECRET_HEADER", "x-build-secret")
	defer os.Unsetenv("BUILDER_WEBHOOK_SECRET_HEADER")

	r := setupRoutes()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/build", nil)
	req.Header.Add("x-build-secret", "foo")

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
