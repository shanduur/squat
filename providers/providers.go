package providers

import "database/sql"

type Describe struct {
	ColumnName      sql.NullString
	ColumnType      sql.NullString
	ColumnLength    sql.NullInt64
	ColumnPrecision sql.NullInt64
	Nullable        sql.NullBool
}

type Provider interface {
	GetProviderName() string
	GetTableDescription(string) ([]Describe, error)
}
