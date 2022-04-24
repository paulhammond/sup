package testutil

import "time"

func Timer() func() time.Duration {
	return func() time.Duration {
		return time.Duration(123 * time.Millisecond)
	}
}
