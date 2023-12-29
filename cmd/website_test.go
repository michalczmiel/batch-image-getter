package cmd

import "testing"

func TestValidateWebsiteCmdArgs(t *testing.T) {
	t.Run("should return error when no arguments are provided", func(t *testing.T) {
		err := validateWebsiteCmdArgs(nil, []string{})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should return error when too many arguments are provided", func(t *testing.T) {
		err := validateWebsiteCmdArgs(nil, []string{
			"https://example.com", "https://example.com",
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should return error when invalid url is provided", func(t *testing.T) {
		err := validateWebsiteCmdArgs(nil, []string{"example.com"})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should not return error when valid url is provided", func(t *testing.T) {
		err := validateWebsiteCmdArgs(nil, []string{"https://example.com"})
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}
