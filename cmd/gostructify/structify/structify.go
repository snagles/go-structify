package structify

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/fatih/structtag"
	"github.com/sirupsen/logrus"
	"github.com/snagles/gostructify/cmd/gostructify/database"
	"github.com/urfave/cli"
)

// Generate is called once per table
func Generate(c *cli.Context, dbName string, t *database.Table) []byte {
	var buf bytes.Buffer
	// create the struct name definition
	fmt.Fprintf(&buf, "// %s is the go struct representation of %s.%s \n type %s struct {\n", structName(t.Name), dbName, t.Name, structName(t.Name))

	// add each column and field based on passed parameters
	for _, column := range t.Columns {
		if c.GlobalString("tags") != "" {
			fmt.Fprintf(&buf, "%s \t %s \t `%s`\n", fieldName(column), fieldType(column, c.GlobalString("nullabletype")), fieldTags(column, strings.Split(c.GlobalString("tags"), ",")))
		} else {
			fmt.Fprintf(&buf, "%s \t %s \t\n", fieldName(column), fieldType(column, c.GlobalString("nullabletype")))
		}
	}

	// add closing struct bracket
	fmt.Fprintf(&buf, "\n}\n\n")

	// add methods
	if c.GlobalString("methods") != "" {
		fmt.Fprintf(&buf, structMethods(structName(t.Name), t.Name, strings.Split(c.GlobalString("methods"), ",")))
	}
	return buf.Bytes()
}

func structName(tablename string) string {
	var runes []rune
	// convert snake case to camel case
	words := strings.Split(tablename, "_")
	for _, w := range words {
		first := true
		for _, c := range w {
			c := c
			if first {
				c = unicode.ToUpper(c)
			}
			runes = append(runes, c)
			first = false
		}
	}
	return string(runes)
}

func fieldName(c database.Column) string {
	var runes []rune
	// convert snake case to camel case
	words := strings.Split(c.Name, "_")
	for _, w := range words {
		first := true
		for _, c := range w {
			c := c
			if first {
				c = unicode.ToUpper(c)
			}
			runes = append(runes, c)
			first = false
		}
	}
	return string(runes)
}

func fieldType(c database.Column, nullabletype string) string {
	// mariadb uses YES, vertica uses t/true
	if c.DatabaseNullable == "YES" || c.DatabaseNullable == "t" || c.DatabaseNullable == "true" {
		switch nullabletype {
		case "guregu":
			return c.Definition.GureguType
		case "sql":
			return c.Definition.SQLType
		}
	}
	return c.Definition.GoType
}

func fieldTags(c database.Column, options []string) string {
	// add things like xml, csv, json, gorm, sqlx tags
	ts := &structtag.Tags{}
	for _, o := range options {
		var t structtag.Tag
		switch o {
		case "json":
			t = JSONTags(c)
		case "xml":
			t = XMLTags(c)
		case "gorm":
			t = GormTags(c)
		case "sqlx":
			t = SQLXTags(c)
		case "csv":
			t = CSVTags(c)
		default:
			logrus.Errorf("Unrecognized tag option: %s", options)
		}
		ts.Set(&t)
	}
	// sort tags according to keys
	sort.Sort(ts)
	return fmt.Sprint(ts)
}

func structMethods(structname string, tablename string, options []string) string {
	var methods []string
	for _, o := range options {
		switch o {
		case "gorm":
			// Add gorm tablename method
			methods = append(methods, GormTableName(structname, tablename))
		default:
			logrus.Errorf("Unrecognized method option: %s", options)
		}
	}
	return strings.Join(methods, "\n")
}
