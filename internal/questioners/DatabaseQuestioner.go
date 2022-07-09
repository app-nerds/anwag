package questioners

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/app-nerds/anwag/internal/answercontext"
	"github.com/app-nerds/anwag/internal/database"
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
				Options: []string{
					database.MongoDB,
					database.MySQL,
					database.Postgres,
				},
			},
		},
		{
			Name: "DbName",
			Prompt: &survey.Input{
				Message: "Name of database",
			},
		},
	}

	if err := survey.Ask(questions, context); err != nil {
		if err == terminal.InterruptErr {
			return ErrInterrupted
		}

		return err
	}

	context.DSN = database.GetDSN(context.WhatTypeOfDatabase, context.DbName)
	context.DBInclude = getDBInclude(context)
	context.DBConnectionCode = getDBConnectionCode(context)

	fmt.Printf("\n\n")
	return nil
}

func getDBInclude(context *answercontext.Context) string {
	switch context.WhatTypeOfDatabase {
	case database.Postgres:
		return `"gorm.io/driver/postgres"
  "gorm.io/gorm"`

	case database.MySQL:
		return `"gorm.io/driver/mysql"
  "gorm.io/gorm"`

	case database.MongoDB:
		return `"github.com/app-nerds/kit/v6/database"`
	}

	return ""
}

func getDBConnectionCode(context *answercontext.Context) string {
	switch context.WhatTypeOfDatabase {
	case database.Postgres:
		return `if db, err = gorm.Open(postgres.Open(config.DSN), &gorm.Config{}); err != nil {
    logger.WithError(err).Fatal("unable to connect to the database")
  }`

	case database.MySQL:
		return `if db, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{}); err != nil {
    logger.WithError(err).Fatal("unable to connect to the database")
  }`

	case database.MongoDB:
		return `if session, err = database.Dial(config.DSN); err != nil {
    logger.WithError(err).Fatal("unable to connect to the database")
  }

  db = session.DB("` + context.DbName + `")
  defer session.Close()`
	}

	return ""
}
