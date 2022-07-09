package generators

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"text/template"

	"github.com/app-nerds/anwag/internal/answercontext"
	"github.com/app-nerds/anwag/internal/errorhandler"
	"github.com/app-nerds/anwag/internal/mappings"
	"github.com/app-nerds/kit/v6/filesystem"
)

func EmptyAppGenerator(context *answercontext.Context, localFS filesystem.FileSystem, templateFS fs.FS) {
	var (
		err       error
		templates *template.Template
		fp        filesystem.WritableFile
	)

	dirs := []string{
		context.AppName + "/cmd/" + context.AppName + "/internal/configuration",
	}

	mapping := []mappings.MappingType{
		{TemplateName: "main.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "main.go")},
		{TemplateName: "Config.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "internal", "configuration", "Config.go")},
		{TemplateName: "go.mod.tmpl", OutputName: "go.mod"},
	}

	for _, dir := range dirs {
		if err = localFS.MkdirAll(dir, 0755); err != nil {
			errorhandler.HandleError(err, fmt.Sprintf("creating base directory `%s`", dir))
		}
	}

	if templates, err = template.ParseFS(templateFS, "templates/emptyapp/*.tmpl"); err != nil {
		errorhandler.HandleError(err, "parsing emptyapp templates")
	}

	for _, m := range mapping {
		if fp, err = localFS.Create(m.Path(context)); err != nil {
			errorhandler.HandleError(err, fmt.Sprintf("creating file `%s`", m.Path(context)))
		}

		defer fp.Close()

		if err = templates.ExecuteTemplate(fp, m.TemplateName, context); err != nil {
			errorhandler.HandleError(err, fmt.Sprintf("executing template %s on file `%s`", m.TemplateName, m.Path(context)))
		}
	}
}
