package parser

import (
	"os"
	"strconv"

	"github.com/askiada/go-sql-test/internal/model"
)

type env struct {
	model.DBCredentials
	SQLFile string
}

func GetEnv() (env, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return env{}, ErrNoDBHost
	}

	dbPortString := os.Getenv("DB_PORT")
	if dbPortString == "" {
		return env{}, ErrNoDBPort
	}

	dbPort, err := strconv.Atoi(dbPortString)
	if err != nil {
		return env{}, ErrDBPortNotNumber
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return env{}, ErrNoDBUser
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return env{}, ErrNoDBPass
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return env{}, ErrNoDBName
	}

	sqlFile := os.Getenv("SQL_FILE")
	if sqlFile == "" {
		return env{}, ErrNoSQLFile
	}

	return env{
		DBCredentials: model.DBCredentials{
			Host: dbHost,
			Port: dbPort,
			User: dbUser,
			Pass: dbPassword,
			Name: dbName,
		},
		SQLFile: sqlFile,
	}, nil
}
