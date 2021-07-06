package providers

import (
	"database/sql"
	"fmt"
)

// ErrNoResult is an error returned when no rows are returned from the query.
var ErrNoResult = fmt.Errorf("table does not exist")

// Describe structure holds definition of the single column in the table.
type Describe struct {
	ColumnName      sql.NullString
	ColumnType      sql.NullString
	ColumnLength    sql.NullInt64
	ColumnPrecision sql.NullInt64
	Nullable        sql.NullBool
}

// Formats struct holds the template formats for the Date, Time and Timestamp
// datatypes.
type Formats struct {
	DateFormat      string `toml:"date"`
	DateTimeFormat  string `toml:"datetime"`
	TimestampFormat string `toml:"timestamp"`
}

// Provider interface is the interface that must be implemented in every
// Database connection package.
type Provider interface {
	ProviderName() string
	GetTableDescription(string) ([]Describe, error)
	DateFormat() string
	DateTimeFormat() string
}
