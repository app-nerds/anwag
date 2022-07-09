package answercontext

import (
	"time"
)

/*
Context houses information used to parse and build templates.
Most of the values come from user input.
*/
type Context struct {
	Year string

	// Basic Questions
	AppName     string
	Title       string
	Description string
	Email       string

	// Type of App Questions
	WhatTypeOfApp string
	WantDatabase  bool

	// Environment Questions
	GithubPath   string
	GithubToken  string
	EnvPrefix    string
	WantFrontend bool

	// Database Questions
	WhatTypeOfDatabase string
	DbName             string
	DSN                string
	DBInclude          string
	DBConnectionCode   string

	GithubSSHPath string
}

/*
NewContext creates a new Context struct
*/
func NewContext() *Context {
	return &Context{
		Year:         time.Now().Format("2006"),
		WantDatabase: false,
	}
}
