package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/app-nerds/postgresr"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v4"
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
	DbName       string
	WantModel    bool

	GithubToken string
}

type tableStruct struct {
	Name    string
	Columns []tableColumn
}

type tableColumn struct {
	Name     string
	DataType string
}

var postgresDatatypes = map[string]string{
	"bigint":                      "int",
	"bit":                         "bool",
	"character varying":           "string",
	"integer":                     "int",
	"numeric":                     "float64",
	"smallint":                    "int",
	"text":                        "string",
	"time without time zone":      "time.Time",
	"timestamp without time zone": "time.Time",
	"uuid":                        "string",
}

func main() {
	var (
		err            error
		fp             *os.File
		sourceFp       fs.File
		templates      *template.Template
		subdirs        []string
		tableNames     []string
		selectedTables []string
		tables         []tableStruct
	)

	values := appValues{
		Year:         time.Now().Format("2006"),
		WantDatabase: false,
	}

	fmt.Printf("\n🙏 Welcome to App Nerds Web Application Generator (%s)\n\n", Version)
	fmt.Printf("Please note that you'll need the following tools to develop in this application:\n")
	fmt.Printf("   Go - https://golang.org\n")
	fmt.Printf("   Swag - https://github.com/swaggo/swag\n")
	fmt.Printf("   NodeJS - https://nodejs.org\n")
	fmt.Printf("   Vue CLI - https://cli.vuejs.org\n")
	fmt.Printf("\n")

	firstQuestions := []*survey.Question{
		{
			Name: "GithubPath",
			Prompt: &survey.Input{
				Message: "Enter the Github path for this application's home repo",
				Help:    "This should be the path to a Github repository. eg. github.com/app-nerds/my-new-app",
			},
			Validate: survey.Required,
		},
		{
			Name:   "CompanyName",
			Prompt: &survey.Input{Message: "Company name", Default: "App Nerds"},
		},
		{
			Name: "Title",
			Prompt: &survey.Input{
				Message: "Title",
				Help:    "This value is used in the browser title and README",
			},
		},
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
			Name: "EnvPrefix",
			Prompt: &survey.Input{
				Message: "Environment variable prefix",
				Help:    "This prefix will be applied to configuration values pulled from the environment. eg. PREFIX_SERVER_HOST. No spaces allowed.",
			},
			Transform: survey.TransformString(func(s string) string {
				upper := strings.ToUpper(s)
				nospaces := strings.ReplaceAll(upper, " ", "")
				notabs := strings.ReplaceAll(nospaces, "\t", "")

				return notabs
			}),
		},
		{
			Name:   "Description",
			Prompt: &survey.Input{Message: "Describe this application in a single sentence"},
		},
		{
			Name:   "Email",
			Prompt: &survey.Input{Message: "Enter your email address"},
		},
		{
			Name: "GithubToken",
			Prompt: &survey.Input{
				Message: "Enter your Github Personal Access Token",
				Help:    "This is used for accessing private repos in Docker builds. See https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token",
			},
		},
		{
			Name:   "WantDatabase",
			Prompt: &survey.Confirm{Message: "Would you like to use a database?", Default: false},
		},
	}

	if err = survey.Ask(firstQuestions, &values); err != nil {
		if err == terminal.InterruptErr {
			fmt.Printf("\nCancelling.\n")
			os.Exit(-1)
		}

		logrus.WithError(err).Fatalf("There was an error!")
	}

	if values.WantDatabase {
		dbQuestions := []*survey.Question{
			{
				Name:     "DbHost",
				Prompt:   &survey.Input{Message: "Database host address", Default: "localhost"},
				Validate: survey.Required,
			},
			{
				Name:     "DbUser",
				Prompt:   &survey.Input{Message: "User name"},
				Validate: survey.Required,
			},
			{
				Name:     "DbPassword",
				Prompt:   &survey.Password{Message: "Password"},
				Validate: survey.Required,
			},
			{
				Name:     "DbName",
				Prompt:   &survey.Input{Message: "Database name"},
				Validate: survey.Required,
			},
			{
				Name:   "WantModel",
				Prompt: &survey.Confirm{Message: "Would you like to generate Go models?"},
			},
		}

		if err = survey.Ask(dbQuestions, &values); err != nil {
			if err == terminal.InterruptErr {
				fmt.Printf("\nCancelling.\n")
				os.Exit(-1)
			}

			logrus.WithError(err).Fatalf("There was an error!")
		}

		if values.WantModel {
			if tableNames, err = getListOfTables(values); err != nil {
				fmt.Printf("Uh oh! There was an error getting a list of tables!\n\n%s\n\nClosing.\n", err.Error())
				os.Exit(-1)
			}

			whichTablesPrompt := &survey.MultiSelect{
				Message: "Which tables should I generate modesl from?",
				Options: tableNames,
			}

			if err = survey.AskOne(whichTablesPrompt, &selectedTables); err != nil {
				if err == terminal.InterruptErr {
					fmt.Printf("\nCancelling.\n")
					os.Exit(-1)
				}

				logrus.WithError(err).Fatalf("There was an error!")
			}

			if tables, err = getColumnsForTables(values, selectedTables); err != nil {
				logrus.WithError(err).Fatalf("Error getting column information!")
			}
		}
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

	fmt.Printf("\nCreating your new application!\n")
	fmt.Printf("   This might take a couple of minutes... ☕️\n")
	fmt.Printf("\n")

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

	/*
	 * Generate models
	 */
	if values.WantDatabase && values.WantModel {
		for _, t := range tables {
			var modelTemplate string
			modelFileName := fmt.Sprintf("%s.go", strcase.ToCamel(t.Name))

			if modelTemplate, err = generateModelFromTable(t); err != nil {
				logrus.WithError(err).Fatalf("error generating table model for %s", t.Name)
			}

			if fp, err = os.Create(modelFileName); err != nil {
				logrus.WithError(err).Fatalf("unable to create model file %s", modelFileName)
			}

			defer fp.Close()

			_, _ = fp.WriteString(modelTemplate)
		}
	}

	/*
	 * Do the steps the initialze the new app
	 */
	if err = modDownload(); err != nil {
		logrus.WithError(err).Fatalf("Error downloading Go modules")
	}

	if err = npmInstall(); err != nil {
		logrus.WithError(err).Fatalf("Error installing Node modules")
	}

	if err = buildNode(); err != nil {
		logrus.WithError(err).Fatalf("Error building NodeJS app")
	}

	fmt.Printf("\n🎉 Congratulations! Your new application is ready.\n")

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

func getListOfTables(values appValues) ([]string, error) {
	// For future, if I support more DBs, we'll switch here
	return getListOfTablesPostgres(values)
}

func getColumnsForTables(values appValues, tableNames []string) ([]tableStruct, error) {
	// For future, if I support more DBs, we'll switch here
	return getColumnsForTablesPostgres(values, tableNames)
}

func getListOfTablesPostgres(values appValues) ([]string, error) {
	var (
		err  error
		db   postgresr.Conn
		rows pgx.Rows

		tablename string
	)

	result := make([]string, 0, 15)

	db, err = postgresr.Connect(context.Background(),
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			values.DbHost, values.DbUser, values.DbPassword,
			values.DbName,
		),
	)

	if err != nil {
		return result, err
	}

	defer db.Close(context.Background())

	sql := `
		SELECT tablename FROM pg_catalog.pg_tables
		WHERE schemaname != 'pg_catalog'
		AND schemaname != 'information_schema'
	`

	if rows, err = db.Query(context.Background(), sql); err != nil {
		return result, err
	}

	for rows.Next() {
		if err = rows.Scan(&tablename); err != nil {
			return result, err
		}

		result = append(result, tablename)
	}

	return result, err
}

func getColumnsForTablesPostgres(values appValues, tableNames []string) ([]tableStruct, error) {
	var (
		err  error
		db   postgresr.Conn
		rows pgx.Rows
	)

	result := make([]tableStruct, 0, 15)

	db, err = postgresr.Connect(context.Background(),
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			values.DbHost, values.DbUser, values.DbPassword,
			values.DbName,
		),
	)

	if err != nil {
		return result, err
	}

	defer db.Close(context.Background())

	for _, tableName := range tableNames {
		t := tableStruct{
			Name:    tableName,
			Columns: make([]tableColumn, 0, 15),
		}

		sql := `
			SELECT 
				column_name, 
				data_type 
			FROM 
				information_schema.columns
			WHERE 
				table_name = $1;
		`

		if rows, err = db.Query(context.Background(), sql, tableName); err != nil {
			return result, err
		}

		var columnname string
		var datatype string

		var actualDataType string
		var ok bool

		for rows.Next() {
			if err = rows.Scan(&columnname, &datatype); err != nil {
				return result, err
			}

			if actualDataType, ok = postgresDatatypes[datatype]; !ok {
				actualDataType = "string"
			}

			c := tableColumn{
				Name:     columnname,
				DataType: actualDataType,
			}

			t.Columns = append(t.Columns, c)
		}

		result = append(result, t)
	}

	return result, nil
}

func generateModelFromTable(table tableStruct) (string, error) {
	var (
		err error
	)

	tmpl := `package main
{{if .HasTime}}
import (
	"time"
)
{{end}}
type {{.Name}} struct {
	{{range .Columns}}{{.Name}} {{.DataType}} ` + "`json:\"{{.JsonName}}\"`" + `
	{{end}}
}
`

	type tempTableColumn struct {
		Name     string
		DataType string
		JsonName string
	}

	data := struct {
		Name    string
		HasTime bool
		Columns []tempTableColumn
	}{
		Name:    strcase.ToCamel(table.Name),
		Columns: make([]tempTableColumn, 0, 15),
	}

	for _, c := range table.Columns {
		newColumn := tempTableColumn{
			Name:     strcase.ToCamel(c.Name),
			DataType: c.DataType,
			JsonName: strcase.ToLowerCamel(c.Name),
		}

		if c.DataType == "time.Time" {
			data.HasTime = true
		}

		data.Columns = append(data.Columns, newColumn)
	}

	t := template.Must(template.New(table.Name).Parse(tmpl))
	b := strings.Builder{}

	if err = t.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

func hasSpace(value string) bool {
	return strings.ContainsAny(value, "  \t")
}

func modDownload() error {
	cmd := exec.Command("go", "mod", "download")
	fmt.Printf("Downloading Go modules...\n")

	err := cmd.Run()
	return err
}

func npmInstall() error {
	var (
		err error
	)

	if err = os.Chdir("app"); err != nil {
		return errors.New("error changing to app directory in npmInstall()")
	}

	cmd := exec.Command("npm", "install")
	fmt.Printf("Installing NodeJS modules...\n")

	if err = cmd.Run(); err != nil {
		return err
	}

	if err = os.Chdir(".."); err != nil {
		return errors.New("error changing back to project root directory in npmInstall()")
	}

	return nil
}

func buildNode() error {
	var (
		err error
	)

	if err = os.Chdir("app"); err != nil {
		return errors.New("error changing to app directory in buildNode()")
	}

	cmd := exec.Command("npm", "run", "build")
	fmt.Printf("Building NodeJS app...\n")

	if err = cmd.Run(); err != nil {
		return err
	}

	if err = os.Chdir(".."); err != nil {
		return errors.New("error changing back to project root directory in buildNode()")
	}

	return nil
}
