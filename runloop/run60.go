package runloop

import (
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// Run the update func until the sentinel shouldStop func returns false,
// or the update func returns shouldContinue=false,
// or the update func returns error.
// Returns the err (if any) from the update func.
func At60Fps(shouldStop func() bool, update func(dt float64) (bool, error)) error {
	var prevTime, dt float64
	tick := time.Tick(time.Second / 60)

	for !shouldStop() {
		// wait for tick
		<-tick

		// calc dt and track time
		now := glfw.GetTime()
		dt = now - prevTime
		prevTime = now
		if dt > 0.02 {
			dt = 0.02
		}

		// Update
		shouldContinue, err := update(dt)

		// Exit on shouldContinue=false or error
		if !shouldContinue || err != nil {
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}
