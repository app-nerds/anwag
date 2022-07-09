package questioners

import "fmt"

var (
	// ErrInterrupted is used when the user cancels this application.
	ErrInterrupted = fmt.Errorf("User interrupted.")
)

