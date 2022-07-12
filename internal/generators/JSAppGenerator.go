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

func JSAppGenerator(context *answercontext.Context, localFS filesystem.FileSystem, templateFS fs.FS) {
	dirs := []string{
		fmt.Sprintf("%s/cmd/%s/internal/configuration", context.AppName, context.AppName),
		fmt.Sprintf("%s/cmd/%s/internal/handlers", context.AppName, context.AppName),
		fmt.Sprintf("%s/cmd/%s/graph/generated", context.AppName, context.AppName),
		fmt.Sprintf("%s/cmd/%s/graph/model", context.AppName, context.AppName),
		fmt.Sprintf("%s/cmd/%s/app/static/js/libraries/nerdwebjs", context.AppName, context.AppName),
		fmt.Sprintf("%s/cmd/%s/app/static/views", context.AppName, context.AppName),
	}

	mapping := []mappings.MappingType{
		{TemplateName: "main.go.tmpl", OutputName: filepath.Join("cmd", context.AppName, "main.go")},
		{TemplateName: "env.tmpl", OutputName: filepath.Join("cmd", context.AppName, ".env")},
		{TemplateName: "env.tmpl", OutputName: filepath.Join("cmd", context.AppName, "env.template")},
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
		{TemplateName: "index.html.tmpl", OutputName: filepath.Join("cmd", context.AppName, "app", "index.html")},
		{TemplateName: "main.js.tmpl", OutputName: filepath.Join("cmd", context.AppName, "app", "main.js")},
		{TemplateName: "manifest.json.tmpl", OutputName: filepath.Join("cmd", context.AppName, "app", "manifest.json")},
		{TemplateName: "nerdwebjs.js.tmpl", OutputName: filepath.Join("cmd", context.AppName, "app", "static", "js", "libraries", "nerdwebjs", "nerdwebjs.js")},
		{TemplateName: "nerdwebjs.js.map.tmpl", OutputName: filepath.Join("cmd", context.AppName, "app", "static", "js", "libraries", "nerdwebjs", "nerdwebjs.js.map")},
		{TemplateName: "nerdwebjs.css.tmpl", OutputName: filepath.Join("cmd", context.AppName, "app", "static", "js", "libraries", "nerdwebjs", "nerdwebjs.css")},
		{TemplateName: "Home.js.tmpl", OutputName: filepath.Join("cmd", context.AppName, "app", "static", "views", "Home.js")},
		{TemplateName: "About.js.tmpl", OutputName: filepath.Join("cmd", context.AppName, "app", "static", "views", "About.js")},
	}

	dir.MakeDirs(localFS, dirs)
	templates.Execute(localFS, templateFS, "templates/jsapp/*.tmpl", mapping, context)
}
