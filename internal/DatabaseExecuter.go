package internal

import (
	"html/template"
	"strings"

	"github.com/iancoleman/strcase"
)

func GenerateModelFromTable(table TableStruct) (string, error) {
	var (
		err error
	)

	tmpl := `package main
{{if .HasTime}}
import (
	"time"
)
{{end}}
type {{.Name}} struct {
	{{range .Columns}}{{.Name}} {{.DataType}} ` + "`json:\"{{.JSONName}}\"`" + `
	{{end}}
}
`

	type tempTableColumn struct {
		Name     string
		DataType string
		JSONName string
	}

	data := struct {
		Name    string
		HasTime bool
		Columns []tempTableColumn
	}{
		Name:    strcase.ToCamel(table.Name),
		Columns: make([]tempTableColumn, 0, 15),
	}

	for _, c := range table.Columns {
		newColumn := tempTableColumn{
			Name:     strcase.ToCamel(c.Name),
			DataType: c.DataType,
			JSONName: strcase.ToLowerCamel(c.Name),
		}

		if c.DataType == "time.Time" {
			data.HasTime = true
		}

		data.Columns = append(data.Columns, newColumn)
	}

	t := template.Must(template.New(table.Name).Parse(tmpl))
	b := strings.Builder{}

	if err = t.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
