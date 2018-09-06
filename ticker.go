package fsutils

import "time"

// CreateTicker it can be either active or not active according
// to the active parameter.
func CreateTicker(intvl time.Duration, active bool) <-chan time.Time {

	// Create the tick chan, the channel will effectively
	// be ignored if verbosity is off.
	var tick <-chan time.Time
	if active {
		tick = time.Tick(500 * time.Millisecond)
	}
	return tick
}
