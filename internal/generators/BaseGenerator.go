package generators

import (
	"fmt"
	"io/fs"
	"text/template"

	"github.com/app-nerds/anwag/internal/answercontext"
	"github.com/app-nerds/anwag/internal/errorhandler"
	"github.com/app-nerds/anwag/internal/mappings"
	"github.com/app-nerds/kit/v6/filesystem"
)

func BaseGenerator(context *answercontext.Context, localFS filesystem.FileSystem, templateFS fs.FS) {
	var (
		err       error
		templates *template.Template
		fp        filesystem.WritableFile
	)

	dirs := []string{
		context.AppName,
	}

	mapping := []mappings.MappingType{
		{TemplateName: "docker-compose.yml.tmpl", OutputName: "docker-compose.yml"},
		{TemplateName: "Dockerfile.tmpl", OutputName: "Dockerfile"},
		{TemplateName: "Makefile.tmpl", OutputName: "Makefile"},
		{TemplateName: "editorconfig.tmpl", OutputName: ".editorconfig"},
		{TemplateName: "gitignore.tmpl", OutputName: ".gitignore"},
		{TemplateName: "env.tmpl", OutputName: ".env"},
		{TemplateName: "env.tmpl", OutputName: "env.template"},
		{TemplateName: "VERSION.tmpl", OutputName: "VERSION"},
		{TemplateName: "README.md.tmpl", OutputName: "README.md"},
		{TemplateName: "CHANGELOG.md.tmpl", OutputName: "CHANGELOG.md"},
	}

	for _, dir := range dirs {
		if err = localFS.MkdirAll(dir, 0755); err != nil {
			errorhandler.HandleError(err, fmt.Sprintf("creating base directory `%s`", dir))
		}
	}

	if templates, err = template.ParseFS(templateFS, "templates/base/*.tmpl"); err != nil {
		errorhandler.HandleError(err, "parsing root templates")
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
