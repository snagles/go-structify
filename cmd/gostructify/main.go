// gostructify is a tool to automate the creation of struct definitions and methods
// representing a table in a database. Given the table name and database connection,
// gostructify will create a new self-contained Go source file containing the
// struct definition and methods using the configured parameters.
//
// The file is created in the same package and directory as the go generate comment.
// This package are largely inspired by https://github.com/golang/tools/tree/master/cmd/stringer
// and https://github.com/Shelnutt2/db2struct

// Typically this process would be run using go generate, like this:
//
// You will then be prompted for host connection details:
// TODO expand usage details
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/sirupsen/logrus"
	"github.com/snagles/gostructify/cmd/gostructify/database"
	"github.com/snagles/gostructify/cmd/gostructify/structify"
	"github.com/urfave/cli"
	"golang.org/x/tools/imports"
)

func main() {
	app := cli.NewApp()

	cli.AppHelpTemplate = tpl

	app.Commands = []cli.Command{
		cli.Command{
			Name: "mariadb",
			Flags: []cli.Flag{
				// username to use
				cli.StringFlag{Name: "username", Usage: "username credentials to use `myuser`"},
				cli.IntFlag{Name: "port", Usage: "database port to connect to `3306`", Value: 3306},
				cli.StringFlag{Name: "hostname", Usage: "hostname to connect to `examplehost.com`"},
				cli.StringFlag{Name: "database", Usage: "database name `application_db`"},
				cli.StringFlag{Name: "tables", Usage: "list of comma separated table names `users,admins`"},
			},
			Usage: "generate structs from a mariadb database",
			Action: func(c *cli.Context) error {
				m := database.MariaDB{Hostname: c.String("hostname"), Username: c.String("username"), Password: getPassword(c, c.String("username")), Port: c.Int("port")}
				process(m, c)
				return nil
			},
		},
		cli.Command{
			Name: "mysql",
			Flags: []cli.Flag{
				// username to use
				cli.StringFlag{Name: "username", Usage: "username credentials to use `myuser`"},
				cli.IntFlag{Name: "port", Usage: "database port to connect to `3306`", Value: 3306},
				cli.StringFlag{Name: "hostname", Usage: "hostname to connect to `examplehost.com`"},
				cli.StringFlag{Name: "database", Usage: "database name `application_db`"},
				cli.StringFlag{Name: "tables", Usage: "list of comma separated table names `users,admins`"},
			},
			Usage: "generate structs from a mysql database",
			Action: func(c *cli.Context) error {
				m := database.MySQL{Hostname: c.String("hostname"), Username: c.String("username"), Password: getPassword(c, c.String("username")), Port: c.Int("port")}
				process(m, c)
				return nil
			},
		},
		cli.Command{
			Name: "postgresql",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "username", Usage: "username credentials to use `myuser`"},
				cli.IntFlag{Name: "port", Usage: "database port to connect to `5432`", Value: 5432},
				cli.StringFlag{Name: "hostname", Usage: "hostname to connect to `examplehost.com`"},
				cli.StringFlag{Name: "database", Usage: "database name `application_db`"},
				cli.StringFlag{Name: "tables", Usage: "list of comma separated table names `users,admins`"},
			},
			Usage: "generate structs from a postgresql database",
			Action: func(c *cli.Context) error {
				p := database.PostgreSQL{Hostname: c.String("hostname"), Username: c.String("username"), Password: getPassword(c, c.String("username")), Port: c.Int("port")}
				process(p, c)
				return nil
			},
		},
		cli.Command{
			Name: "vertica",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "dsn", Usage: "dsn name `vertica`"},
				cli.StringFlag{Name: "database", Usage: "database name `application_db`"},
				cli.StringFlag{Name: "tables", Usage: "list of comma separated table names `users,admins`"},
			},
			Usage: "generate structs from a vertica database",
			Action: func(c *cli.Context) error {
				v := database.Vertica{DSN: c.String("dsn")}
				process(v, c)
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		// directorie or file to parse
		cli.StringFlag{
			Name:  "directory",
			Usage: "directory to search for generate flags `.`",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "file to search for generate flags `.`",
		},
		// field types and annotations
		cli.StringFlag{
			Name:  "nullabletype",
			Usage: "preferred handling of nullable types `sql,guregu`",
		},
		cli.StringFlag{
			Name:  "tags",
			Usage: "list of comma delimited tag options `json,gorm,sqlx,xml,csv`",
		},
		cli.StringFlag{
			Name:  "methods",
			Usage: "list of comma delmited method options `gorm`",
		},
		cli.BoolFlag{
			Name:  "stdin",
			Usage: "use stdin for passing of password",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "password for database connection",
			EnvVar: "GOSTRUCTIFY_PASSWORD",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "dry-run prints output to screen instead of a file",
		},
	}

	app.Name = "gostructify"
	app.HelpName = "gostructify"
	app.Version = "0.0.1"
	app.UsageText = "Generate go structs, types, and methods from database tables"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "snagles",
		},
	}
	app.Run(os.Args)
}

func process(d database.Database, c *cli.Context) {
	g := Generator{}
	if c.GlobalString("directory") != "" && isDirectory(c.GlobalString("directory")) {
		g.parsePackageDir(c.GlobalString("directory"))
	} else {
		g.parsePackageFiles(strings.Split(c.GlobalString("file"), ","))
	}

	for _, db := range strings.Split(c.String("database"), ",") {
		// build the struct for each table
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "// Code generated by gostructify command \"%s\"; DO NOT EDIT.\n", strings.Join(os.Args, " "))
		fmt.Fprintf(&buf, "package %s", g.pkg.name)
		fmt.Fprintf(&buf, "\n")
		fmt.Fprintf(&buf, "\n")

		for _, table := range strings.Split(c.String("tables"), ",") {
			t, err := d.Build(db, table)
			if err != nil {
				logrus.Fatalf("Failed to generate file for %s.%s: %s", db, table, err)
			}
			buf.Write(structify.Generate(c, db, t))
		}

		// Write the file unformatted
		if c.GlobalBool("dry-run") {
			f, err := ioutil.TempFile("", "tmp")
			if err != nil {
				logrus.Fatalf("failed to write temp file for dry run: %s", err)
			}

			defer os.Remove(f.Name())

			_, err = f.Write(buf.Bytes())
			if err != nil {
				logrus.Fatalf("writing tempfile output: %s", err)
			}

			// Run go imports to gofmt and goimports the file
			src, err := imports.Process(f.Name(), buf.Bytes(), nil)
			if err != nil {
				logrus.Fatalf("processing imports: %s", err)
			}
			fmt.Println(string(src))
		} else {
			output := c.GlobalString("file")
			if output == "" {
				baseName := fmt.Sprintf("%s_gostructify.go", db)
				output = filepath.Join(g.pkg.dir, strings.ToLower(baseName))
			}
			err := ioutil.WriteFile(output, buf.Bytes(), 0644)
			if err != nil {
				logrus.Fatalf("writing output: %s", err)
			}

			// Run go imports to gofmt and goimports the file
			src, err := imports.Process(output, buf.Bytes(), nil)
			if err != nil {
				logrus.Fatalf("processing imports: %s", err)
			}

			err = ioutil.WriteFile(output, src, 0644)
			if err != nil {
				logrus.Fatalf("writing output: %s", err)
			}
			fmt.Printf("Wrote file: %s", output)
		}
	}
}

func getPassword(c *cli.Context, username string) string {
	if c.GlobalString("password") == "" {
		if c.GlobalBool("stdin") {
			fmt.Printf("Enter password for user %s (if no password, leave blank): ", username)
			password, err := gopass.GetPasswd()
			if err != nil {
				logrus.Fatalf("Failed to retrieve password: %s", err)
			}
			return string(password)
		}
		logrus.Fatalf("Missing password. Password either needs to be set as an env variable `GOSTRUCTIFY_PASSWORD=password`, passed on the command line, or using the --stdin option to be prompted")
	}

	return c.GlobalString("password")
}

const tpl = `NAME:
   {{.HelpName}} (v{{.Version}}) - {{.UsageText}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   gostructify --directory ~/go/src/github.com/snagles/testproject --tags gorm,sqlx,json --methods gorm mariadb --hostname 127.0.0.1 --username root --tables all_data_types --database test
   {{if len .Authors}}
AUTHOR: {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
   {{end}} `
