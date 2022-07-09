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
)

func main() {
	var (
		err error
		// basicQuestioner       internal.Questioner
		// appSettingsQuestioner internal.Questioner
		// dbQuestioner          internal.Questioner
		// fp                    filesystem.WritableFile
		// sourceFp              fs.File
		// templates             *template.Template
		// subdirs               []string
		// tableNames            []string
		// selectedTables        []string
		// tables                []tableStruct
	)

	localFS := localfs.NewLocalFS()
	context := answercontext.NewContext()

	// basicQuestioner = internal.BasicsQuestioner{}
	// appSettingsQuestioner = internal.AppSettingsQuestioner{}
	// dbQuestioner = internal.DatabaseQuestioner{}

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
	case typesofapps.EmptyApp:
		generators.EmptyAppGenerator(context, localFS, emptyAppFs)

	case typesofapps.ApiApp:
		generators.ApiAppGenerator(context, localFS, apiAppFs)
	}

	// rootFsMapping := []internal.MappingType{
	// 	{TemplateName: "main.go.tmpl", OutputName: "main.go"},
	// 	{TemplateName: "ClientApp.go.tmpl", OutputName: "ClientApp.go", IsFrontend: true},
	// 	{TemplateName: "Config.go.tmpl", OutputName: "Config.go"},
	// 	{TemplateName: "docker-compose.yml.tmpl", OutputName: "docker-compose.yml"},
	// 	{TemplateName: "Dockerfile.tmpl", OutputName: "Dockerfile"},
	// 	{TemplateName: "go.mod.tmpl", OutputName: "go.mod"},
	// 	{TemplateName: "Makefile.tmpl", OutputName: "Makefile"},
	// 	{TemplateName: "VersionController.go.tmpl", OutputName: "VersionController.go"},
	// 	{TemplateName: "VersionController_test.go.tmpl", OutputName: "VersionController_test.go"},
	// 	{TemplateName: "editorconfig.tmpl", OutputName: ".editorconfig"},
	// 	{TemplateName: "gitignore.tmpl", OutputName: ".gitignore"},
	// 	{TemplateName: "create.sql.tmpl", OutputName: "devops/sql/create.sql", IsDatabase: true},
	// 	{TemplateName: "env.tmpl", OutputName: ".env"},
	// 	{TemplateName: "VERSION.tmpl", OutputName: "VERSION"},
	// 	{TemplateName: "README.md.tmpl", OutputName: "README.md"},
	// 	{TemplateName: "CHANGELOG.md.tmpl", OutputName: "CHANGELOG.md"},
	// 	{TemplateName: "go.yml.tmpl", OutputName: ".github/workflows/go.yml"},
	// }

	// appFsMapping := map[string]string{
	// 	"browserlistrc.tmpl":       "app/.browserlistrc",
	// 	"env.tmpl":                 "app/.env",
	// 	"env.production.tmpl":      "app/.env.production",
	// 	"eslintrc.js.tmpl":         "app/.eslintrc.js",
	// 	"gitignore.tmpl":           "app/.gitignore",
	// 	"prettierrc.tmpl":          "app/.prettierrc",
	// 	"babel.config.js.tmpl":     "app/babel.config.js",
	// 	"jest.config.js.tmpl":      "app/jest.config.js",
	// 	"package.json.tmpl":        "app/package.json",
	// 	"index.html.tmpl":          "app/public/index.html",
	// 	"index2.html.tmpl":         "app/dist/index.html",
	// 	"styles.css.tmpl":          "app/src/assets/css/styles.css",
	// 	"AppFooter.vue.tmpl":       "app/src/modules/ui/components/AppFooter.vue",
	// 	"AppHeader.vue.tmpl":       "app/src/modules/ui/components/AppHeader.vue",
	// 	"Logo.vue.tmpl":            "app/src/modules/ui/components/Logo.vue",
	// 	"AlertService.js.tmpl":     "app/src/modules/ui/services/AlertService.js",
	// 	"VersionService.js.tmpl":   "app/src/modules/version/services/VersionService.js",
	// 	"index.js.tmpl":            "app/src/router/index.js",
	// 	"Home.vue.tmpl":            "app/src/views/Home.vue",
	// 	"PageTwo.vue.tmpl":         "app/src/views/PageTwo.vue",
	// 	"App.vue.tmpl":             "app/src/App.vue",
	// 	"BaseURL.js.tmpl":          "app/src/BaseURL.js",
	// 	"HttpInterceptors.js.tmpl": "app/src/HttpInterceptors.js",
	// 	"main.js.tmpl":             "app/src/main.js",
	// 	"example.spec.js.tmpl":     "app/tests/unit/example.spec.js",
	// }

	// staticFilesMapping := map[string]string{
	// 	"templates/app/favicon.ico":        "app/public/favicon.ico",
	// 	"templates/app/app-nerds-logo.jpg": "app/src/assets/images/app-nerds-logo.jpg",
	// }

	// fmt.Printf("\nCreating your new application!\n")
	// fmt.Printf("   This might take a couple of minutes... ☕️\n")
	// fmt.Printf("\n")

	// /*
	//  * Create directories
	//  */
	// if err = localFS.Mkdir(context.AppName, 0755); err != nil {
	// 	logrus.WithError(err).Errorf("error attempting to create application directory '%s'", context.AppName)
	// }

	// if err = localFS.Chdir(context.AppName); err != nil {
	// 	logrus.WithError(err).Errorf("error changing to new directory '%s'", context.AppName)
	// }

	// subdirs = []string{
	// 	"docs",
	// 	".github/workflows",
	// }

	// if context.WantFrontend {
	// 	subdirs = append(subdirs,
	// 		"app",
	// 		"app/dist/css",
	// 		"app/dist/js",
	// 		"app/dist/img",
	// 		"app/public",
	// 		"app/src/assets",
	// 		"app/src/assets/css",
	// 		"app/src/assets/images",
	// 		"app/src/modules",
	// 		"app/src/modules/ui",
	// 		"app/src/modules/ui/components",
	// 		"app/src/modules/ui/services",
	// 		"app/src/modules/version",
	// 		"app/src/modules/version/services",
	// 		"app/src/router",
	// 		"app/src/views",
	// 		"app/tests/unit",
	// 	)
	// }

	// if context.WantDatabase {
	// 	subdirs = append(subdirs, "devops/sql")
	// }

	// for _, dirname := range subdirs {
	// 	if err = localFS.MkdirAll(dirname, 0755); err != nil {
	// 		logrus.WithError(err).Fatalf("unable to create subdirectory %s", dirname)
	// 	}
	// }

	// /*
	//  * Process root filesystem templates
	//  */
	// if templates, err = template.ParseFS(rootFs, "templates/*.tmpl"); err != nil {
	// 	logrus.WithError(err).Fatal("error parsing root templates files")
	// }

	// for _, mt := range rootFsMapping {
	// 	if !context.WantFrontend && mt.IsFrontend {
	// 		continue
	// 	}

	// 	if !context.WantDatabase && mt.IsDatabase {
	// 		continue
	// 	}

	// 	if fp, err = os.Create(mt.OutputName); err != nil {
	// 		logrus.WithError(err).Fatalf("unable to create file %s for writing", mt.OutputName)
	// 	}

	// 	defer fp.Close()

	// 	if err = templates.ExecuteTemplate(fp, mt.TemplateName, context); err != nil {
	// 		logrus.WithError(err).Fatalf("unable to exeucte template %s", mt.TemplateName)
	// 	}
	// }

	// /*
	//  * Process app filesystem templates
	//  */
	// if context.WantFrontend {
	// 	if templates, err = template.ParseFS(rootFs, "templates/app/*.tmpl"); err != nil {
	// 		logrus.WithError(err).Fatal("error parsing app template files")
	// 	}

	// 	for templateName, outputPath := range appFsMapping {
	// 		if fp, err = localFS.Create(outputPath); err != nil {
	// 			logrus.WithError(err).Fatalf("unable to create file %s for writing", outputPath)
	// 		}

	// 		defer fp.Close()

	// 		if err = templates.ExecuteTemplate(fp, templateName, context); err != nil {
	// 			logrus.WithError(err).Fatalf("unable to execute template %s", templateName)
	// 		}
	// 	}

	// 	/*
	// 	 * Copy static assets
	// 	 */
	// 	for sourceFile, targetFile := range staticFilesMapping {
	// 		if sourceFp, err = rootFs.Open(sourceFile); err != nil {
	// 			logrus.WithError(err).Fatalf("unable to read static asset %s", sourceFile)
	// 		}

	// 		defer sourceFp.Close()

	// 		if fp, err = localFS.Create(targetFile); err != nil {
	// 			logrus.WithError(err).Fatalf("unable to create file %s for writing", targetFile)
	// 		}

	// 		defer fp.Close()

	// 		if _, err = io.Copy(fp, sourceFp); err != nil {
	// 			logrus.WithError(err).Fatalf("error copying static asset %s to %s", sourceFile, targetFile)
	// 		}
	// 	}
	// }

	// /*
	//  * Generate models
	//  */
	// if context.WantDatabase && context.WantModel {
	// 	for _, t := range context.Tables {
	// 		var modelTemplate string
	// 		modelFileName := fmt.Sprintf("%s.go", strcase.ToCamel(t.Name))

	// 		if modelTemplate, err = internal.GenerateModelFromTable(t); err != nil {
	// 			logrus.WithError(err).Fatalf("error generating table model for %s", t.Name)
	// 		}

	// 		if fp, err = localFS.Create(modelFileName); err != nil {
	// 			logrus.WithError(err).Fatalf("unable to create model file %s", modelFileName)
	// 		}

	// 		defer fp.Close()

	// 		_, _ = fp.WriteString(modelTemplate)
	// 	}
	// }

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
