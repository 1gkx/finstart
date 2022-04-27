package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var dbDriver = "postgres"

type MigrateCongig struct {
	Folder  string
	Version int
	Log     log.Logger
}

func newMigrateConfig() *MigrateCongig {
	return &MigrateCongig{
		Folder:  "",
		Version: -1,
		Log:     log.NewLogfmtLogger(os.Stdout),
	}
}

func Migrate(dsn string, c *MigrateCongig) error {

	mc := newMigrateConfig()
	if c.Log != nil {
		mc.Log = c.Log
	}
	if c.Version > 0 {
		mc.Version = c.Version
	}
	if len(c.Folder) > 0 {
		mc.Folder = c.Folder
	}

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		mc.Log.Log("event", "error", "desc", err)
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		mc.Log.Log("event", "error", "desc", err)
		return err
	}

	folder := fmt.Sprintf("file://%s", mc.Folder)
	m, err := migrate.NewWithDatabaseInstance(folder, dbDriver, driver)
	if err != nil {
		mc.Log.Log("event", "error", "desc", err)
		return err
	}

	// m.Force(mc.Version)

	if err := m.Migrate(uint(mc.Version)); err != nil && !errors.Is(migrate.ErrNoChange, err) {
		if errors.Is(migrate.ErrNoChange, err) {
			mc.Log.Log("event", "migrate", "desc", err)
			return nil
		}
		return err
	}
	return nil
}
