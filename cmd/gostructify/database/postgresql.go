package database

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq" // postgres driver
)

// PostgreSQL and its methods builds the common table and column structure for formatting
type PostgreSQL struct {
	Hostname string
	Port     int
	Username string
	Password string
}

const (
	postgresColumnQuery = "SELECT column_name, data_type, is_nullable FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = '%s' AND table_name = '%s'"
)

// Build retrieves the table and column schema information for the table and database name to be generated
// it then parses and converts the database specific types into the correct go types
func (p PostgreSQL) Build(database, table string) (*Table, error) {
	db, err := sql.Open("postgres", p.connectionString(database))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf(postgresColumnQuery, database, table))
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

		// remove postgres sizes
		r := regexp.MustCompile("\\([0-9]+\\)")
		c.DatabaseType = r.ReplaceAllString(c.DatabaseType, "")
		r = regexp.MustCompile("\\([0-9]+,[0-9]+\\)")
		c.DatabaseType = r.ReplaceAllString(c.DatabaseType, "")

		var ok bool
		if c.Definition, ok = postgresTypeMap[c.DatabaseType]; !ok {
			return nil, fmt.Errorf("Unrecognized column type field: %s", c.DatabaseType)
		}
		t.Columns = append(t.Columns, c)
	}
	if len(t.Columns) == 0 {
		return nil, fmt.Errorf("No rows returned from the information schema for TABLE_SCHEMA %s and TABLE_NAME %s", database, table)
	}
	return &t, nil
}

func (p *PostgreSQL) connectionString(database string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", p.Username, p.Password, database)
}

var postgresTypeMap = map[string]ColumnDefinition{
	// integer fields
	"bigint":      ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"bigserial":   ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"integer":     ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"int":         ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"int4":        ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"smallint":    ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"int2":        ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"smallserial": ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"serial2":     ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"serial":      ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"serial4":     ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	// string fields
	"char":              ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"character":         ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"character varying": ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"varchar":           ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"cidr":              ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"inet":              ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"json":              ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"macaddr":           ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"text":              ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"xml":               ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	// time fields
	"date": ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"time": ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"time without time zone":      ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"time with time zone":         ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"timetz":                      ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"timestamp":                   ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"timestamp without time zone": ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"timestamp with time zone":    ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"timestamptz":                 ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"interval":                    ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	// float fields
	"double precision": ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"float8":           ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"money":            ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"numeric":          ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"decimal":          ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"real":             ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	// binary fields
	"bit":         ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"bit varying": ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"bytea":       ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"jsonb":       ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	// bool
	"boolean": ColumnDefinition{GoType: "bool", GureguType: "null.Bool", SQLType: "null.Bool"},
}
