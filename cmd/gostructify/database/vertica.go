package database

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/alexbrainman/odbc" // odbc connection driver for vertica
)

// Vertica and its methods builds the common table and column structure for formatting
type Vertica struct {
	DSN string
}

const (
	verticaColumnQuery = "SELECT column_name, data_type, is_nullable FROM v_catalog.columns WHERE table_schema = ? AND table_name = ?"
)

// Build retrieves the table and column schema information for the table and database name to be generated
// it then parses and converts the database specific types into the correct go types
func (v Vertica) Build(database, table string) (*Table, error) {
	db, err := sql.Open("odbc", v.connectionString(database))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(verticaColumnQuery, database, table)
	if err != nil {
		return nil, err
	}

	t := Table{Name: table}
	defer rows.Close()
	for rows.Next() {
		c := Column{}

		err := rows.Scan(&c.Name, &c.DatabaseType, &c.DatabaseNullable)
		if err != nil {
			return nil, err
		}

		// remove vertica sizes
		r := regexp.MustCompile("\\([0-9]+\\)")
		c.DatabaseType = r.ReplaceAllString(c.DatabaseType, "")
		r = regexp.MustCompile("\\([0-9]+,[0-9]+\\)")
		c.DatabaseType = r.ReplaceAllString(c.DatabaseType, "")

		var ok bool
		if c.Definition, ok = verticaTypeMap[c.DatabaseType]; !ok {
			return nil, fmt.Errorf("Unrecognized column type field: %s", c.DatabaseType)
		}
		t.Columns = append(t.Columns, c)
	}
	if len(t.Columns) == 0 {
		return nil, fmt.Errorf("No rows returned from the information schema for TABLE_SCHEMA %s and TABLE_NAME %s", database, table)
	}
	return &t, nil
}

func (v *Vertica) connectionString(database string) string {
	// ex: "Driver={odbc};Servername=example;Database=example;Port=5433;uid=dbadmin;pwd=test;ResultBufferSize=0;ConnectionLoadBalance=1;"
	return fmt.Sprintf("DSN=%s;ResultBufferSize=0;ConnectionLoadBalance=1;", v.DSN)
}

var verticaTypeMap = map[string]ColumnDefinition{
	// integer fields
	"tinyint":   ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"int":       ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"smallint":  ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"mediumint": ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"bigint":    ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	// string fields
	"char":              ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"character":         ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"varchar":           ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"character varying": ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"long varchar":      ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	// time fields
	"date":          ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"datetime":      ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"smalldatetime": ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"time":          ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"timestamp":     ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	// float fields
	"float":            ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"float8":           ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"float16":          ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"float32":          ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"float64":          ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"numeric":          ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"number":           ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"double precision": ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"real precision":   ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	// binary fields
	"binary":         ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"varbinary":      ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"binary varying": ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"bytea":          ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"raw":            ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	// bool
	"boolean": ColumnDefinition{GoType: "bool", GureguType: "null.Bool", SQLType: "null.Bool"},
}
