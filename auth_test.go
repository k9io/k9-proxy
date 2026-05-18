package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthenticate_API_MissingKey(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	Authenticate_API()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthenticate_API_InvalidFormat(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("API_KEY", "badformat")

	Authenticate_API()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthenticate_API_ValidKey(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("API_KEY", "company123:apikey456")

	Authenticate_API()(c)

	if c.GetString("company_uuid") != "company123" {
		t.Fatalf("expected company_uuid=company123, got %q", c.GetString("company_uuid"))
	}
	if c.GetString("api_key") != "apikey456" {
		t.Fatalf("expected api_key=apikey456, got %q", c.GetString("api_key"))
	}
}

func TestCacheAuthenticate_API(t *testing.T) {
	dir := t.TempDir()
	Config = &Configuration{}
	Config.Proxy.Cache_Dir = dir

	apiKey := "company123:secret"
	cacheFile := filepath.Join(dir, "api.cache")
	Write_Cache(cacheFile, apiKey)

	if !Cache_Authenticate_API(apiKey) {
		t.Fatal("expected true for matching key")
	}
	if Cache_Authenticate_API("company123:wrong") {
		t.Fatal("expected false for wrong key")
	}
}
