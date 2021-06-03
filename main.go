package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
)

var (
	Version string = "development"

	//go:embed templates/*.tmpl
	rootFs embed.FS

	//go:embed templates/app/*.tmpl
	//go:embed templates/app/favicon.ico
	//go:embed templates/app/app-nerds-logo.jpg
	appFs embed.FS
)

type appValues struct {
	Year string

	CompanyName string
	Title       string
	AppName     string
	EnvPrefix   string
	Description string
	Email       string
	GithubPath  string

	WantDatabase bool
	DbHost       string
	DbUser       string
	DbPassword   string

	GithubToken string
}

func main() {
	var (
		err       error
		fp        *os.File
		sourceFp  fs.File
		templates *template.Template
		subdirs   []string
	)

	values := appValues{
		Year:         time.Now().Format("2006"),
		WantDatabase: false,
	}

	charsReplacer := strings.NewReplacer("-", "", "_", "", " ", "")

	fmt.Printf("Welcome to App Nerds Web Application Generator (%s)\n\n", Version)

	values.GithubPath = stringPrompt("Enter the Github path for this application's home repo", "")
	lastPartOfRepo := values.GithubPath[strings.LastIndex(values.GithubPath, "/")+1:]
	strippedLastPartOfRepo := charsReplacer.Replace(lastPartOfRepo)
	capsLastPartOfRepo := strings.ToUpper(strippedLastPartOfRepo)

	values.CompanyName = stringPrompt("Company name", "")
	values.Title = stringPrompt("Title (used in browser title and header)", "")
	values.AppName = stringPrompt("Application name (no spaces)", strippedLastPartOfRepo)
	values.EnvPrefix = stringPrompt("Environment variable prefix (no spaces, all caps)", capsLastPartOfRepo)
	values.Description = stringPrompt("Describe this application in a single sentance", "")
	values.Email = stringPrompt("Enter your email address", "")
	values.GithubToken = stringPrompt("Enter your Github Personal Access Token", "")
	values.WantDatabase = yesNoPrompt("Would you like to use a database?")

	if values.WantDatabase {
		values.DbHost = stringPrompt("Database host address", "localhost")
		values.DbUser = stringPrompt("User name", strippedLastPartOfRepo)
		values.DbPassword = stringPrompt("Password", "password")
	}

	rootFsMapping := map[string]string{
		"main.go.tmpl":                   "main.go",
		"ClientApp.go.tmpl":              "ClientApp.go",
		"Config.go.tmpl":                 "Config.go",
		"docker-compose.yml.tmpl":        "docker-compose.yml",
		"Dockerfile.tmpl":                "Dockerfile",
		"go.mod.tmpl":                    "go.mod",
		"Makefile.tmpl":                  "Makefile",
		"VersionController.go.tmpl":      "VersionController.go",
		"VersionController_test.go.tmpl": "VersionController_test.go",
		"editorconfig.tmpl":              ".editorconfig",
		"gitignore.tmpl":                 ".gitignore",
		"create.sql.tmpl":                "devops/sql/create.sql",
		"env.tmpl":                       ".env",
		"VERSION.tmpl":                   "VERSION",
		"README.md.tmpl":                 "README.md",
		"CHANGELOG.md.tmpl":              "CHANGELOG.md",
		"go.yml.tmpl":                    ".github/workflows/go.yml",
	}

	appFsMapping := map[string]string{
		"browserlistrc.tmpl":       "app/.browserlistrc",
		"env.tmpl":                 "app/.env",
		"env.production.tmpl":      "app/.env.production",
		"eslintrc.js.tmpl":         "app/.eslintrc.js",
		"gitignore.tmpl":           "app/.gitignore",
		"prettierrc.tmpl":          "app/.prettierrc",
		"babel.config.js.tmpl":     "app/babel.config.js",
		"jest.config.js.tmpl":      "app/jest.config.js",
		"package.json.tmpl":        "app/package.json",
		"index.html.tmpl":          "app/public/index.html",
		"index2.html.tmpl":         "app/dist/index.html",
		"styles.css.tmpl":          "app/src/assets/css/styles.css",
		"AppFooter.vue.tmpl":       "app/src/modules/ui/components/AppFooter.vue",
		"AppHeader.vue.tmpl":       "app/src/modules/ui/components/AppHeader.vue",
		"Logo.vue.tmpl":            "app/src/modules/ui/components/Logo.vue",
		"AlertService.js.tmpl":     "app/src/modules/ui/services/AlertService.js",
		"VersionService.js.tmpl":   "app/src/modules/version/services/VersionService.js",
		"index.js.tmpl":            "app/src/router/index.js",
		"Home.vue.tmpl":            "app/src/views/Home.vue",
		"PageTwo.vue.tmpl":         "app/src/views/PageTwo.vue",
		"App.vue.tmpl":             "app/src/App.vue",
		"BaseURL.js.tmpl":          "app/src/BaseURL.js",
		"HttpInterceptors.js.tmpl": "app/src/HttpInterceptors.js",
		"main.js.tmpl":             "app/src/main.js",
		"example.spec.js.tmpl":     "app/tests/unit/example.spec.js",
	}

	staticFilesMapping := map[string]string{
		"templates/app/favicon.ico":        "app/public/favicon.ico",
		"templates/app/app-nerds-logo.jpg": "app/src/assets/images/app-nerds-logo.jpg",
	}

	/*
	 * Create directories
	 */
	if err = os.Mkdir(values.AppName, 0755); err != nil {
		logrus.WithError(err).Errorf("error attempting to create application directory '%s'", values.AppName)
	}

	if err = os.Chdir(values.AppName); err != nil {
		logrus.WithError(err).Errorf("error changing to new directory '%s'", values.AppName)
	}

	subdirs = []string{
		"app",
		"app/dist/css",
		"app/dist/js",
		"app/dist/img",
		"app/public",
		"app/src/assets",
		"app/src/assets/css",
		"app/src/assets/images",
		"app/src/modules",
		"app/src/modules/ui",
		"app/src/modules/ui/components",
		"app/src/modules/ui/services",
		"app/src/modules/version",
		"app/src/modules/version/services",
		"app/src/router",
		"app/src/views",
		"app/tests/unit",
		"devops/sql",
		"docs",
		".github/workflows",
	}

	for _, dirname := range subdirs {
		if err = os.MkdirAll(dirname, 0755); err != nil {
			logrus.WithError(err).Fatalf("unable to create subdirectory %s", dirname)
		}
	}

	/*
	 * Process root filesystem templates
	 */
	if templates, err = template.ParseFS(rootFs, "templates/*.tmpl"); err != nil {
		logrus.WithError(err).Fatal("error parsing root templates files")
	}

	for templateName, outputPath := range rootFsMapping {
		if fp, err = os.Create(outputPath); err != nil {
			logrus.WithError(err).Fatalf("unable to create file %s for writing", outputPath)
		}

		defer fp.Close()

		if err = templates.ExecuteTemplate(fp, templateName, values); err != nil {
			logrus.WithError(err).Fatalf("unable to exeucte template %s", templateName)
		}
	}

	/*
	 * Process app filesystem templates
	 */
	if templates, err = template.ParseFS(appFs, "templates/app/*.tmpl"); err != nil {
		logrus.WithError(err).Fatal("error parsing app template files")
	}

	for templateName, outputPath := range appFsMapping {
		if fp, err = os.Create(outputPath); err != nil {
			logrus.WithError(err).Fatalf("unable to create file %s for writing", outputPath)
		}

		defer fp.Close()

		if err = templates.ExecuteTemplate(fp, templateName, values); err != nil {
			logrus.WithError(err).Fatalf("unable to execute template %s", templateName)
		}
	}

	/*
	 * Copy static assets
	 */
	for sourceFile, targetFile := range staticFilesMapping {
		if sourceFp, err = appFs.Open(sourceFile); err != nil {
			logrus.WithError(err).Fatalf("unable to read static asset %s", sourceFile)
		}

		defer sourceFp.Close()

		if fp, err = os.Create(targetFile); err != nil {
			logrus.WithError(err).Fatalf("unable to create file %s for writing", targetFile)
		}

		defer fp.Close()

		if _, err = io.Copy(fp, sourceFp); err != nil {
			logrus.WithError(err).Fatalf("error copying static asset %s to %s", sourceFile, targetFile)
		}
	}

	fmt.Printf("\n🎉 Congratulations! Your new application is ready.\n")
	fmt.Printf("Please note that you'll need the following tools to develop in this application:\n")
	fmt.Printf("   Go - https://golang.org\n")
	fmt.Printf("   Swag - https://github.com/swaggo/swag\n")
	fmt.Printf("   NodeJS - https://nodejs.org\n")
	fmt.Printf("   Vue CLI - https://cli.vuejs.org\n")
	fmt.Printf("\nTo begin execute the following:\n\n")
	fmt.Printf("   cd %s/app\n", values.AppName)
	fmt.Printf("   npm install && npm run build\n")
	fmt.Printf("   cd ..\n")
	fmt.Printf("   go get\n")

	if values.WantDatabase {
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
}

func stringPrompt(label, defaultValue string) string {
	var (
		err    error
		result string
	)

	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}

	if result, err = prompt.Run(); err != nil {
		logrus.WithError(err).Fatalf("error asking for '%s'", label)
	}

	return result
}

func yesNoPrompt(label string) bool {
	var (
		err       error
		selection string
	)

	prompt := promptui.Select{
		Label: label,
		Items: []string{"No", "Yes"},
	}

	if _, selection, err = prompt.Run(); err != nil {
		logrus.WithError(err).Fatalf("error asking for '%s'", label)
	}

	if selection == "Yes" {
		return true
	}

	return false
}
