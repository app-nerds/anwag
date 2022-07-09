package questioners

import (
  "fmt"

  "github.com/app-nerds/anwag/internal/answercontext"
  "github.com/app-nerds/anwag/internal/database"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func AskDatabaseQuestions(context *answercontext.Context) error {
  if !context.WantDatabase {
    return nil
  }

  fmt.Printf("DATABASE QUESTIONS\n")
  fmt.Printf("---------------------\n")

	questions := []*survey.Question{
    {
      Name: "WhatTypeOfDatabase",
      Prompt: &survey.Select{
        Message: "What type of database are you using?",
        Options: []string {
          database.MongoDB,
          database.Postgres,
        },
      },
    },
    {
      Name: "DbHost",
      Prompt: &survey.Input{
        Message: "Database host address", 
        Default: "localhost",
      },
    },
    {
      Name: "DbName",
      Prompt: &survey.Input{
        Message: "Name of database",
      },
    },
    {
      Name: "DbUser",
      Prompt: &survey.Input{
        Message: "User name used to connect to this database",
      },
    },
    {
      Name: "DbPassword",
      Prompt: &survey.Password{
        Message: "Password used to connect to this database",
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
