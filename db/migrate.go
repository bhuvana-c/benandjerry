package db

import (
	"fmt"
	"strings"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
	"github.com/zalora/benandjerry/config"
)

func RunMigrateScripts(dbConnConfig *config.PgConn, migrationScriptsPath string) error {
	allErrors, ok := migrate.UpSync(dbConnConfig.GetURL(), migrationScriptsPath)
	if !ok {
		errString := []string{"Error running migration scripts : \n"}
		for _, err := range allErrors {
			errString = append(errString, err.Error())
		}
		return fmt.Errorf(strings.Join(errString, "\n"))
	}
	return nil
}
