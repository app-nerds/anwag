package questioners

import (
  "fmt"
  "strings"

  "github.com/app-nerds/anwag/internal/answercontext"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func AskBasicQuestions(context *answercontext.Context) error {
  fmt.Printf("BASIC QUESTIONS\n")
  fmt.Printf("---------------\n")

	questions := []*survey.Question{
		{
			Name: "AppName",
			Prompt: &survey.Input{
				Message: "Application name",
				Help:    "This is used primary as the final executable name. There should be no spaces in this name. Spaces will be replaced with hyphens",
			},
			Transform: survey.TransformString(func(s string) string {
				lower := strings.ToLower(s)
				nospaces := strings.ReplaceAll(lower, " ", "-")
				notabs := strings.ReplaceAll(nospaces, "\t", "-")

				return notabs
			}),
		},
		{
			Name: "Title",
			Prompt: &survey.Input{
				Message: "Title",
				Help:    "This value is used in the browser title and README",
			},
		},
		{
			Name:   "Description",
			Prompt: &survey.Input{Message: "Describe this application in a single sentence"},
		},
		{
			Name:   "Email",
			Prompt: &survey.Input{Message: "Enter your email address"},
		},
	}

	if err := survey.Ask(questions, context); err != nil {
		if err == terminal.InterruptErr {
			return ErrInterrupted
		}

		return err
	}

  fmt.Printf("\n\n")
  return nil
}
