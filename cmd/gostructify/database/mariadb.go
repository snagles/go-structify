package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // mysql driver for mariadb support
)

// MariaDB and its methods builds the common table and column structure for formatting
type MariaDB struct {
	Hostname string
	Port     int
	Username string
	Password string
}

const (
	mariaDBColumnQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ?"
)

// Build retrieves the table and column schema information for the table and database name to be generated
// it then parses and converts the database specific types into the correct go types
func (m MariaDB) Build(database, table string) (*Table, error) {
	db, err := sql.Open("mysql", m.connectionString(database))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(mariaDBColumnQuery, database, table)
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
		var ok bool
		if c.Definition, ok = mariaDBTypeMap[c.DatabaseType]; !ok {
			return nil, fmt.Errorf("Unrecognized column type field: %s", c.DatabaseType)
		}
		t.Columns = append(t.Columns, c)
	}
	if len(t.Columns) == 0 {
		return nil, fmt.Errorf("No rows returned from the information schema for TABLE_SCHEMA %s and TABLE_NAME %s", database, table)
	}
	return &t, nil
}

func (m *MariaDB) connectionString(database string) string {
	// ex: "username:password@tcp(hostname:hostport)/databasename?parseTime=True"
	return fmt.Sprintf("%s:%v@tcp(%s:%v)/%s?parseTime=true", m.Username, m.Password, m.Hostname, m.Port, database)
}

var mariaDBTypeMap = map[string]ColumnDefinition{
	// integer fields
	"tinyint":   ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"int":       ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"smallint":  ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"mediumint": ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	"bigint":    ColumnDefinition{GoType: "int", GureguType: "null.Int", SQLType: "sql.NullInt64"},
	// string fields
	"char":       ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"set":        ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"enum":       ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"varchar":    ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"tinytext":   ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"longtext":   ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"mediumtext": ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	"text":       ColumnDefinition{GoType: "string", GureguType: "null.String", SQLType: "sql.NullString"},
	// time fields
	"date":      ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"datetime":  ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"time":      ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"year":      ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	"timestamp": ColumnDefinition{GoType: "time.Time", GureguType: "null.Time", SQLType: "mysql.NullTime"},
	// float fields
	"decimal": ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"double":  ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	"float":   ColumnDefinition{GoType: "float64", GureguType: "null.Float", SQLType: "sql.NullFloat64"},
	// binary fields
	"binary":     ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"blob":       ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"tinyblob":   ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"longblob":   ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"mediumblob": ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
	"varbinary":  ColumnDefinition{GoType: "[]byte", GureguType: "[]byte", SQLType: "[]byte"},
}
