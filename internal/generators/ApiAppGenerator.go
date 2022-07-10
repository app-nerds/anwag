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
		fmt.Sprintf("%s/cmd/%s/internal/handlers", context.AppName, context.AppName),
		fmt.Sprintf("%s/cmd/%s/graph/generated", context.AppName, context.AppName),
		fmt.Sprintf("%s/cmd/%s/graph/model", context.AppName, context.AppName),
	}

	mapping := []mappings.MappingType{
		{TemplateName: "main.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "main.go")},
		{TemplateName: "Config.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "internal", "configuration", "Config.go")},
		{TemplateName: "go.mod.tmpl", OutputName: "go.mod"},
		{TemplateName: "VersionHandler.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "internal", "handlers", "VersionHandler.go")},
		{TemplateName: "tools.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "tools.go")},
		{TemplateName: "gqlgen.yml.tmpl", OutputName: filepath.Join("cmd", context.AppName, "gqlgen.yml")},
		{TemplateName: "gqlgen.yml.tmpl", OutputName: filepath.Join("cmd", context.AppName, "gqlgen.yml")},
		{TemplateName: "generated.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "graph", "generated", "generated.go")},
		{TemplateName: "models_gen.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "graph", "model", "models_gen.go")},
		{TemplateName: "resolver.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "graph", "resolver.go")},
		{TemplateName: "schema.graphqls.tmpl", OutputName: filepath.Join("cmd", context.AppName, "graph", "schema.graphqls")},
		{TemplateName: "schema.resolvers.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "graph", "schema.resolvers.go")},
	}

	dir.MakeDirs(localFS, dirs)
	templates.Execute(localFS, templateFS, "templates/apiapp/*.tmpl", mapping, context)
}