package config

import (
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type db struct {
	Host   string
	Port   string
	DbName string
	User   string
	Pass   string
}

type migration struct {
	Folder  string
	Version int
}

type Config struct {
	HttpPort  string
	Migration *migration
	DB        *db
}

func Init() *Config {

	_ = godotenv.Load()

	version, _ := strconv.Atoi(os.Getenv("VERSION"))

	return &Config{
		HttpPort: os.Getenv("PORT"),
		Migration: &migration{
			Folder:  os.Getenv("MIGRATION_FOLDER"),
			Version: version,
		},
		DB: &db{
			Host:   os.Getenv("POSTGRES_HOST"),
			Port:   os.Getenv("POSTGRES_PORT"),
			DbName: os.Getenv("POSTGRES_DB"),
			User:   os.Getenv("POSTGRES_USER"),
			Pass:   os.Getenv("POSTGRES_PASSWORD"),
		},
	}
}

func (c *Config) DbDsn() string {

	const (
		cDataBaseURLParameter = "database"
		cSSLModeParameter     = "sslmode"
		cDebugURLParameter    = "debug"
	)

	connParams := make(url.Values)
	connParams.Set(cDataBaseURLParameter, c.DB.DbName)
	// connParams.Add(cDebugURLParameter, "false")
	connParams.Add(cSSLModeParameter, "disable")

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.DB.User, c.DB.Pass),
		Host:     net.JoinHostPort(c.DB.Host, c.DB.Port),
		RawQuery: connParams.Encode(),
	}).String()
}
