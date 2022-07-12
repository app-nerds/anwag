package main

import (
	"embed"
	"errors"
	"fmt"
	// "io"
	// "io/fs"
	"os"
	"os/exec"
	"strings"
	// "text/template"

	// "github.com/app-nerds/anwag/internal"
	"github.com/app-nerds/anwag/internal/answercontext"
	"github.com/app-nerds/anwag/internal/generators"
	"github.com/app-nerds/anwag/internal/questioners"
	"github.com/app-nerds/anwag/internal/typesofapps"
	"github.com/app-nerds/kit/v6/filesystem"
	"github.com/app-nerds/kit/v6/filesystem/localfs"
	// "github.com/iancoleman/strcase"
	// "github.com/sirupsen/logrus"
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
	fmt.Printf("   Swag - https://github.com/swaggo/swag\n")
	fmt.Printf("   NodeJS - https://nodejs.org\n")
	fmt.Printf("   Vue CLI - https://cli.vuejs.org\n")
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

	// Generate!
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

	// /*
	//  * Do the steps the initialze the new app
	//  */
	// if context.WantFrontend {
	// 	if err = npmInstall(localFS); err != nil {
	// 		logrus.WithError(err).Fatalf("Error installing Node modules")
	// 	}

	// 	if err = buildNode(localFS); err != nil {
	// 		logrus.WithError(err).Fatalf("Error building NodeJS app")
	// 	}
	// }

	// if err = modDownload(); err != nil {
	// 	logrus.WithError(err).Fatalf("Error downloading Go modules")
	// }

	// if err = gitInit(); err != nil {
	// 	logrus.WithError(err).Fatalf("Error initializing Git repository")
	// }

	fmt.Printf("\n🎉 Congratulations! Your new application is ready.\n")

	if context.WantDatabase {
		fmt.Printf("\n👩‍💻 Since you opted to have a database setup, in a new terminal window run: \n")
		fmt.Printf("   make start-postgres\n")
		fmt.Printf("\nThis will start a Postgres database in a Docker container. ")
		fmt.Printf("Now, in your first terminal window, run the following to start the app:\n")
	}

	fmt.Printf("   make run\n\n")

	fmt.Printf("\n")
	fmt.Printf("☁ Please note that this project is also setup to use Github Actions.\n")
	fmt.Printf("Anytime a tag in the format of 'v*.*.*' is pushed to origin\n")
	fmt.Printf("the code is built and a Release is created. Ensure you update\n")
	fmt.Printf("the VERSION and CHANGELOG.md files before pushing a tag.\n")

	fmt.Printf("\n")
	fmt.Printf("Before this works, however, you'll need to setup a secret in your\n")
	fmt.Printf("repository called ACTIONS_TOKEN. This should contain a Personal\n")
	fmt.Printf("Access Token that has access to private repositories.\n")

	fmt.Printf("\n\n%+v\n\n", context)
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

func modDownload() error {
	var (
		err error
	)

	fmt.Printf("Downloading Go modules...\n")
	cmd := exec.Command("go", "mod", "download")

	if err = cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("go", "get")
	err = cmd.Run()

	return err
}

func npmInstall(localFS filesystem.FileSystem) error {
	var (
		err error
	)

	if err = localFS.Chdir("app"); err != nil {
		return errors.New("error changing to app directory in npmInstall()")
	}

	cmd := exec.Command("npm", "install")
	fmt.Printf("Installing NodeJS modules...\n")

	if err = cmd.Run(); err != nil {
		return err
	}

	if err = localFS.Chdir(".."); err != nil {
		return errors.New("error changing back to project root directory in npmInstall()")
	}

	return nil
}

func buildNode(localFS filesystem.FileSystem) error {
	var (
		err error
	)

	if err = localFS.Chdir("app"); err != nil {
		return errors.New("error changing to app directory in buildNode()")
	}

	cmd := exec.Command("npm", "run", "build")
	fmt.Printf("Building NodeJS app...\n")

	if err = cmd.Run(); err != nil {
		return err
	}

	if err = localFS.Chdir(".."); err != nil {
		return errors.New("error changing back to project root directory in buildNode()")
	}

	return nil
}

func gitInit() error {
	var (
		err error
	)

	fmt.Printf("Initialize Git repository... \n")
	cmd := exec.Command("git", "init")

	err = cmd.Run()
	return err
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
