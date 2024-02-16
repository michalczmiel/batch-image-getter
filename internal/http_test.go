package internal

import "testing"

func TestGetRandomUserAgent(t *testing.T) {
	t.Run("should return a random user agent", func(t *testing.T) {
		userAgent := getRandomUserAgent()

		if userAgent == "" {
			t.Error("Expected user agent, got empty string")
		}
	})
}
