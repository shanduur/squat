package providers

import (
	"database/sql"
	"fmt"
)

var ErrNoResult = fmt.Errorf("table does not exist")

type Describe struct {
	ColumnName      sql.NullString
	ColumnType      sql.NullString
	ColumnLength    sql.NullInt64
	ColumnPrecision sql.NullInt64
	Nullable        sql.NullBool
}

type Formats struct {
	DateFormat      string `toml:"date"`
	DateTimeFormat  string `toml:"datetime"`
	TimestampFormat string `toml:"timestamp"`
}

type Provider interface {
	ProviderName() string
	GetTableDescription(string) ([]Describe, error)
	DateFormat() string
	DateTimeFormat() string
}
