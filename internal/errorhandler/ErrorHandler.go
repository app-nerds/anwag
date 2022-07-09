package errorhandler

import (
  "fmt"
  "os"
)

func HandleError(err error, action string) {
  fmt.Printf("An unexpected error occured while performing '%s':\n", action)
  fmt.Printf("  %s\n\n", err.Error())
  os.Exit(-1)
}
