package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/app-nerds/anwag/internal/answercontext"
	"github.com/app-nerds/anwag/internal/generators"
	"github.com/app-nerds/anwag/internal/questioners"
	"github.com/app-nerds/anwag/internal/typesofapps"
	"github.com/app-nerds/kit/v6/filesystem"
	"github.com/app-nerds/kit/v6/filesystem/localfs"
)

var (
	// Version is populated when building. This is the version of this application
	Version string = "development"

	//go:embed templates/base/*.tmpl
	baseFs embed.FS

	//go:embed templates/emptyapp/*.tmpl
	emptyAppFs embed.FS

	//go:embed templates/apiapp/*.tmpl
	apiAppFs embed.FS

	//go:embed templates/jsapp/*.tmpl
	jsAppFs embed.FS

	//go:embed templates/cronapp/*.tmpl
	cronAppFs embed.FS
)

func main() {
	var (
		err error
	)

	localFS := localfs.NewLocalFS()
	context := answercontext.NewContext()

	fmt.Printf("\n🙏 Welcome to App Nerds Web Application Generator (%s)\n\n", Version)
	fmt.Printf("Please note that you'll need the following tools to develop in this application:\n")
	fmt.Printf("   Go - https://golang.org\n")
	fmt.Printf("   make - Comes with most Linux and MacOS distros")
	fmt.Printf("   Docker - https://www.docker.com\n")
	fmt.Printf("\n")

	if err = questioners.AskBasicQuestions(context); err != nil {
		exitOnInterrupt(err)
	}

	if err = questioners.AskTypeOfAppQuestions(context); err != nil {
		exitOnInterrupt(err)
	}

	if err = questioners.AskEnvQuestions(context); err != nil {
		exitOnInterrupt(err)
	}

	if err = questioners.AskDatabaseQuestions(context); err != nil {
		exitOnInterrupt(err)
	}

	context.GithubSSHPath = generateSSHURL(context.GithubPath)

	/*
	 * Generate
	 */
	generators.BaseGenerator(context, localFS, baseFs)

	switch context.WhatTypeOfApp {
	case typesofapps.ApiApp:
		generators.ApiAppGenerator(context, localFS, apiAppFs)

	case typesofapps.JSApp:
		generators.JSAppGenerator(context, localFS, jsAppFs)

	case typesofapps.CronApp:
		generators.CronAppGenerator(context, localFS, cronAppFs)

	default:
		generators.EmptyAppGenerator(context, localFS, emptyAppFs)
	}

	/*
	 * Do the steps the initialze the new app
	 */
	if err = makeSetup(localFS, context); err != nil {
		fmt.Printf("Unexpected error: %s\n\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("\n🎉 Congratulations! Your new application is ready.\n")

	if context.WantDatabase {
		fmt.Printf("\n👩‍💻 Since you opted to have a database setup, in a new terminal window run: \n")
		fmt.Printf("   make start-database\n")
		fmt.Printf("\nThis will start database server in a Docker container.\n")
	}

	fmt.Printf("\nGo ahead and run your application by executing the following:\n")
	fmt.Printf("   make run\n\n")

	os.Chdir(context.AppName)
}

func hasSpace(value string) bool {
	return strings.ContainsAny(value, "  \t")
}

func generateSSHURL(githubPath string) string {
	var (
		split []string
	)

	split = strings.Split(githubPath, "/")

	if len(split) < 3 {
		return fmt.Sprintf("https://%s", githubPath)
	}

	return fmt.Sprintf("git@%s:%s/%s.git", split[0], split[1], split[2])
}

func makeSetup(localFS filesystem.FileSystem, context *answercontext.Context) error {
	var err error

	fmt.Printf("Initializing application...\n")

	if err = localFS.Chdir(context.AppName); err != nil {
		return fmt.Errorf("error changing directory to %s: %w", context.AppName, err)
	}

	cmd := exec.Command("make", "setup")

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("error running 'make setup': %w", err)
	}

	return nil
}

func exitOnInterrupt(err error) {
	if errors.Is(err, questioners.ErrInterrupted) {
		fmt.Printf("\nCancelling.\n")
		os.Exit(-1)
	}

	fmt.Printf("\nThere was an unexpected error.\n")
	fmt.Printf("   %s\n", err.Error())
	os.Exit(-1)
}
