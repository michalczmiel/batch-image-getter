package internal

import (
	"testing"
)

func TestSleeperDuration(t *testing.T) {
	t.Run("return duration when min and max are the same", func(t *testing.T) {
		min := 5
		max := 5
		sleeper := &realSleeper{MinInterval: min, MaxInterval: max}
		result := sleeper.duration()

		if result != min {
			t.Errorf("want %d got %d", min, result)
		}
	})

	t.Run("return random duration when min and max are different", func(t *testing.T) {
		min := 1
		max := 5
		sleeper := &realSleeper{MinInterval: min, MaxInterval: max}

		results := make([]int, 0, 25)
		for i := 0; i < 25; i++ {
			results = append(results, sleeper.duration())
		}

		for _, r := range results {
			if r < min || r > max {
				t.Errorf("want a number between %d and %d got %d", min, max, r)
			}
		}
	})
}
