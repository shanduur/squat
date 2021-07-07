package postgres

import (

	// embed is used here for including describe.sql file during compilation.
	_ "embed"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/shanduur/squat/config"
	"github.com/shanduur/squat/providers"
)

type pgConfig struct {
	ProviderName string            `toml:"provider-name"`
	Address      string            `toml:"address"`
	Port         uint16            `toml:"port"`
	Database     string            `toml:"database"`
	User         string            `toml:"user"`
	Password     string            `toml:"password"`
	Formats      providers.Formats `toml:"formats"`
}

//go:embed describe.sql
var describeQuery string

var dbURL = "postgresql://%v:%v/%v?user=%v&password=%v"

type PgProvider struct {
	cfg pgConfig
}

func New(configPath string) (PgProvider, error) {
	var pg PgProvider

	err := pg.Initialize(configPath)
	if err != nil {
		return pg, fmt.Errorf("unable to initialize: %s", err.Error())
	}

	return pg, nil
}

// Initialize reads config and returns the provider with configuration read
// from the config. By default it is called by New function, but can be used standalone.
func (pg *PgProvider) Initialize(configPath string) (err error) {
	err = config.ReadTOML(&pg.cfg, configPath)
	if err != nil {
		return fmt.Errorf("unable to read config: %s", err.Error())
	}

	return nil
}

func (pg PgProvider) ProviderName() string {
	return pg.cfg.ProviderName
}

func (pg PgProvider) GetTableDescription(name string) (dsc []providers.Describe, err error) {
	conn, err := connect(pg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %s", err.Error())
	}

	rows, err := conn.Query(describeQuery, name)
	if err != nil {
		return nil, fmt.Errorf("unable to execute statement: %s", err.Error())
	}
	defer rows.Close()

	d := providers.Describe{}
	for rows.Next() {
		err = rows.Scan(&d.ColumnName, &d.ColumnType, &d.ColumnLength, &d.ColumnPrecision, &d.Nullable)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %s", err.Error())
		}
		dsc = append(dsc, d)
	}

	if len(dsc) < 1 {
		err = providers.ErrNoResult
	}

	return
}

func (pg PgProvider) DateFormat() string {
	return pg.cfg.Formats.DateFormat
}

func (pg PgProvider) DateTimeFormat() string {
	return pg.cfg.Formats.DateTimeFormat
}

func connect(pg PgProvider) (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(pgx.ConnConfig{
		Host:     pg.cfg.Address,
		Port:     pg.cfg.Port,
		Database: pg.cfg.Database,
		User:     pg.cfg.User,
		Password: pg.cfg.Password,
	})
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %s", err.Error())
		return
	}

	return
}
