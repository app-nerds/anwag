package generators

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/app-nerds/anwag/internal/answercontext"
	"github.com/app-nerds/anwag/internal/dir"
	"github.com/app-nerds/anwag/internal/mappings"
	"github.com/app-nerds/anwag/internal/templates"
	"github.com/app-nerds/kit/v6/filesystem"
)

func ApiAppGenerator(context *answercontext.Context, localFS filesystem.FileSystem, templateFS fs.FS) {
	dirs := []string{
		fmt.Sprintf("%s/cmd/%s/internal/configuration", context.AppName, context.AppName),
	}

	mapping := []mappings.MappingType{
		{TemplateName: "main.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "main.go")},
		{TemplateName: "Config.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "internal", "configuration", "Config.go")},
		{TemplateName: "go.mod.tmpl", OutputName: "go.mod"},
	}

	dir.MakeDirs(localFS, dirs)
	templates.Execute(localFS, templateFS, "templates/emptyapp/*.tmpl", mapping, context)
}
