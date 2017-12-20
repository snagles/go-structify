package database

type (
	// Table contains all column definitions
	Table struct {
		Columns []Column
		Name    string
	}

	// Column contains the necessary information to generate the column struct field
	Column struct {
		Name string
		// Database representation
		DatabaseType     string
		DatabaseNullable string
		Definition       ColumnDefinition
	}

	// ColumnDefinition contains the necessary information for the struct field type
	ColumnDefinition struct {
		// Used for creation of the struct
		GoType     string
		GureguType string
		SQLType    string
	}
)
