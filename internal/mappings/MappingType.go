package mappings

import (
  "path/filepath"

  "github.com/app-nerds/anwag/internal/answercontext"
)

type MappingType struct {
	TemplateName string
	OutputName   string
	IsFrontend   bool
	IsDatabase   bool
}

func (mt MappingType) Path(context *answercontext.Context) string {
  return filepath.Join(context.AppName, mt.OutputName)
}
