package questioners

import (
  "fmt"
  "os"

  "github.com/app-nerds/anwag/internal/answercontext"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func AskEnvQuestions(context *answercontext.Context) error {
  fmt.Printf("ENVIRONMENT QUESTIONS\n")
  fmt.Printf("---------------------\n")

	questions := []*survey.Question{
		{
			Name: "GithubPath",
			Prompt: &survey.Input{
				Message: "Enter the Github path for this application's home repo",
				Help:    "This should be the path to a Github repository. eg. github.com/app-nerds/my-new-app",
				Suggest: func(toComplete string) []string {
					return []string{
						"github.com/app-nerds/" + context.AppName,
					}
				},
			},
			Validate: survey.Required,
		},
		{
			Name: "GithubToken",
			Prompt: &survey.Input{
				Message: "Enter your Github Personal Access Token",
				Help:    "This is used for accessing private repos in Docker builds. See https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token",
				Suggest: func(toComplete string) []string {
					value := os.Getenv("GITHUB_TOKEN")

					if value == "" {
						return []string{}
					}

					return []string{
						value,
					}
				},
			},
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

