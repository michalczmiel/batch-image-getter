package cmd

import (
	"testing"
)

func TestValidateFileCmdArguments(t *testing.T) {
	t.Run("should return error when no arguments are provided", func(t *testing.T) {
		err := validateFileCmdArguments([]string{})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should return error when too many arguments are provided", func(t *testing.T) {
		err := validateFileCmdArguments([]string{
			"images1.txt", "images2.txt",
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should not return error when valid url is provided", func(t *testing.T) {
		err := validateFileCmdArguments([]string{"images.txt"})
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}
