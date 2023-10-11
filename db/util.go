package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultDbName = "items"
)

type DatabaseConnConf struct {
	Username string
	Password string
	Hostname string
	Port     string
	Database string
}

func ReadDatabaseConnectionConfig() *DatabaseConnConf {
	dbConf := &DatabaseConnConf{}
	dbConf.Username = os.Getenv("DB_USERNAME")
	dbConf.Password = os.Getenv("DB_PASSWORD")
	dbConf.Hostname = os.Getenv("DB_HOSTNAME")
	dbConf.Port = os.Getenv("DB_PORT")
	dbConf.Database = defaultDbName
	return dbConf
}

func dbOpenConnection(dbConf *DatabaseConnConf) (*sql.DB, error) {
	dbConn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConf.Username, dbConf.Password, dbConf.Hostname, dbConf.Port, dbConf.Database))
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func Ping(dbConf *DatabaseConnConf) error {
	dbConn, err := dbOpenConnection(dbConf)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	return dbConn.Ping()
}

func CreateDatabase(dbConf *DatabaseConnConf) error {
	dbConn, err := dbOpenConnection(dbConf)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	sqlQuery := fmt.Sprintf("CREATE DATABASE %s", dbConf.Database)
	_, err = dbConn.Exec(sqlQuery)
	if err != nil {
		return err
	}
	return nil
}

func EnsureTable(dbConf *DatabaseConnConf, tableName string) error {
	dbConn, err := dbOpenConnection(dbConf)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	sqlQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, num FLOAT)", tableName)
	_, err = dbConn.Exec(sqlQuery)
	if err != nil {
		return err
	}
	return nil
}
