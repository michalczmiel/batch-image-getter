package cmd

import (
	"testing"
)

func TestValidateArguments(t *testing.T) {
	t.Run("should return error when no arguments are provided", func(t *testing.T) {
		err := validateArguments([]string{})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should return error when too many arguments are provided", func(t *testing.T) {
		err := validateArguments([]string{
			"https://example.com", "https://example.com",
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should return error when invalid url is provided", func(t *testing.T) {
		err := validateArguments([]string{"example.com"})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should not return error when valid url is provided", func(t *testing.T) {
		err := validateArguments([]string{"https://example.com"})
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}
