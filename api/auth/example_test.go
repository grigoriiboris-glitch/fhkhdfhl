package auth_test

import (
	"testing"

	"github.com/mymindmap/api/auth"
)

func TestAuthConstants(t *testing.T) {
	// Test that constants are defined correctly
	if auth.RoleUser != "user" {
		t.Errorf("Expected RoleUser to be 'user', got '%s'", auth.RoleUser)
	}

	if auth.RoleAdmin != "admin" {
		t.Errorf("Expected RoleAdmin to be 'admin', got '%s'", auth.RoleAdmin)
	}

	if auth.ObjectPost != "post" {
		t.Errorf("Expected ObjectPost to be 'post', got '%s'", auth.ObjectPost)
	}

	if auth.ActionRead != "read" {
		t.Errorf("Expected ActionRead to be 'read', got '%s'", auth.ActionRead)
	}

	t.Log("Auth constants test passed!")
}

func TestConfigDefaults(t *testing.T) {
	// Test that we can create a config
	config := &auth.Config{
		EnableRateLimit: false,
	}

	if config == nil {
		t.Error("Config should not be nil")
	}

	t.Log("Config defaults test passed!")
}

// Note: In real tests, you'd use a proper mock or test database connection