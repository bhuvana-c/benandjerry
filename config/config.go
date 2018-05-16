package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type PgConn struct {
	Host     string
	Database string
	User     string
	Password string
	Port     uint16
	SslMode  string
}

func (pgConn *PgConn) GetURL() string {
	sslMode := ""
	if pgConn.SslMode != "" {
		sslMode = "?sslmode=" + pgConn.SslMode
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", pgConn.User,
		pgConn.Password, pgConn.Host, strconv.Itoa(int(pgConn.Port)), pgConn.Database, sslMode)
}

func (pgConn *PgConn) Decode(dbUrl string) error {
	url, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	if url.User == nil {
		return fmt.Errorf("username and password not provided")
	}
	password, _ := url.User.Password()
	hostAndPort := strings.SplitN(url.Host, ":", 2)
	if len(hostAndPort) != 2 {
		return fmt.Errorf("unable to parse hostname and port")
	}
	dbPort, err := strconv.ParseUint(hostAndPort[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port %s", hostAndPort[1])
	}
	port := uint16(dbPort)
	queryStrings := url.Query()
	*pgConn = PgConn{
		Host:     hostAndPort[0],
		Database: strings.TrimLeft(url.Path, "/"),
		User:     url.User.Username(),
		Password: password,
		Port:     port,
		SslMode:  queryStrings.Get("sslmode"),
	}
	return nil
}

type Config struct {
	LogLevel             string `envconfig:"LOG_LEVEL"`
	LogFormat            string `envconfig:"LOG_FORMAT"`
	ListenPort           uint16 `envconfig:"LISTEN_PORT" required:"true"`
	PostgresConn         PgConn `envconfig:"POSTGRES_URL" required:"true"`
	MigrationScriptsPath string `envconfig:"MIGRATION_SCRIPTS_PATH" required:"true"`
}

func Load() (*Config, error) {
	var config Config
	err := envconfig.Process("BAJ", &config)
	return &config, err
}
