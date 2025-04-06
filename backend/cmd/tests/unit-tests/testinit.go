package apis

import (
	"os"
	"testing"
)

func TestGetEnvWithDefault(t *testing.T) {
	// Test case 1: Environment variable is set
	os.Setenv("TEST_KEY", "test_value")
	result := getEnvWithDefault("TEST_KEY", "default_value")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test case 2: Environment variable is not set
	os.Unsetenv("TEST_KEY")
	result = getEnvWithDefault("TEST_KEY", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnvWithDefault_EmptyKey(t *testing.T) {
	// Test case 3: Empty key should return default value
	result := getEnvWithDefault("", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}
