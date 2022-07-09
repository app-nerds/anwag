package templates

import (
	"fmt"
	"io/fs"
	gotemplates "text/template"

	"github.com/app-nerds/anwag/internal/answercontext"
	"github.com/app-nerds/anwag/internal/errorhandler"
	"github.com/app-nerds/anwag/internal/mappings"
	"github.com/app-nerds/kit/v6/filesystem"
)

func Execute(localFS filesystem.FileSystem, templateFS fs.FS, templatePath string, mappings []mappings.MappingType, context *answercontext.Context) {
	var (
		err error
		t   *gotemplates.Template
		fp  filesystem.WritableFile
	)

	if t, err = gotemplates.ParseFS(templateFS, templatePath); err != nil {
		errorhandler.HandleError(err, "parsing emptyapp templates")
	}

	for _, m := range mappings {
		if fp, err = localFS.Create(m.Path(context)); err != nil {
			errorhandler.HandleError(err, fmt.Sprintf("creating file `%s`", m.Path(context)))
		}

		defer fp.Close()

		if err = t.ExecuteTemplate(fp, m.TemplateName, context); err != nil {
			errorhandler.HandleError(err, fmt.Sprintf("executing template %s on file `%s`", m.TemplateName, m.Path(context)))
		}
	}
}
