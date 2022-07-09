package questioners

import (
  "fmt"

  "github.com/app-nerds/anwag/internal/answercontext"
  "github.com/app-nerds/anwag/internal/typesofapps"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func AskTypeOfAppQuestions(context *answercontext.Context) error {
  fmt.Printf("TYPE OF APP QUESTIONS\n")
  fmt.Printf("---------------------\n")

	questions := []*survey.Question{
		{
			Name: "WhatTypeOfApp",
			Prompt: &survey.Select{
				Message: "What type of application are you making?",
        Options: []string{
          typesofapps.EmptyApp,
          typesofapps.ApiApp,
          typesofapps.VueApp,
          typesofapps.CronApp,
        },
			},
		},
		{
			Name:   "WantDatabase",
			Prompt: &survey.Confirm{Message: "Would you like to use a database?", Default: false},
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

