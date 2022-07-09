package dir

import (
	"fmt"

	"github.com/app-nerds/anwag/internal/errorhandler"
	"github.com/app-nerds/kit/v6/filesystem"
)

func MakeDirs(localFS filesystem.FileSystem, dirs []string) {
	var (
		err error
	)

	for _, dir := range dirs {
		if err = localFS.MkdirAll(dir, 0755); err != nil {
			errorhandler.HandleError(err, fmt.Sprintf("creating directory '%s'", dir))
		}
	}
}
