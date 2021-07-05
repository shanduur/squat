package informix

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/alexbrainman/odbc"
	"github.com/shanduur/squat/config"
	"github.com/shanduur/squat/providers"
)

type ifxConfig struct {
	ProviderName   string            `toml:"provider-name"`
	DataSourceName string            `toml:"odbc-data-source"`
	Formats        providers.Formats `toml:"formats"`
}

//go:embed describe.sql
var describeQuery string

type IfxProvider struct {
	cfg ifxConfig
}

func New(configPath string) (IfxProvider, error) {
	var ifx IfxProvider

	err := ifx.Initialize(configPath)
	if err != nil {
		return ifx, fmt.Errorf("unable to initialize: %s", err.Error())
	}

	return ifx, nil
}

func (ifx *IfxProvider) Initialize(configPath string) (err error) {
	err = config.ReadTOML(&ifx.cfg, configPath)
	if err != nil {
		return fmt.Errorf("unable to read config: %s", err.Error())
	}

	return nil
}

func (ifx IfxProvider) GetTableDescription(name string) (dsc []providers.Describe, err error) {
	conn, err := connect(ifx)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %s", err.Error())
	}

	stmt, err := conn.Prepare(describeQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %s", err.Error())
	}

	rows, err := stmt.Query(name)
	if err != nil {
		return nil, fmt.Errorf("unable to execute statement: %s", err.Error())
	}
	defer rows.Close()

	d := providers.Describe{}
	for rows.Next() {
		rows.Scan(&d.ColumnName, &d.ColumnType, &d.ColumnLength, &d.ColumnPrecision, &d.Nullable)
		dsc = append(dsc, d)
	}

	return
}

func (ifx IfxProvider) ProviderName() string {
	return ifx.cfg.ProviderName
}

func (ifx IfxProvider) DateFormat() string {
	return ifx.cfg.Formats.DateFormat
}

func (ifx IfxProvider) DateTimeFormat() string {
	return ifx.cfg.Formats.DateTimeFormat
}

func connect(ifx IfxProvider) (conn *sql.DB, err error) {
	conn, err = sql.Open("odbc", fmt.Sprintf("DSN=%s", ifx.cfg.DataSourceName))
	if err != nil {
		return nil, fmt.Errorf("openning connection failed: %s", err.Error())
	}

	return
}
