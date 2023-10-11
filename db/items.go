package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

const (
	ItemsTableName = "items"
)

type ItemsDbClient interface {
	SelectAllItems() ([]Item, error)
	GetItemById(id string) (*Item, error)
	GetRandomItem() (*Item, error)
	InsertItem(item Item) (string, error)
	DeleteItemById(id string) (bool, error)
}

type itemsDbClientImpl struct {
	databaseConnConf *DatabaseConnConf
	tableName        string
}

var _ ItemsDbClient = (*itemsDbClientImpl)(nil)

func NewItemsDbClient(databaseConnConf *DatabaseConnConf) ItemsDbClient {
	return &itemsDbClientImpl{
		databaseConnConf: databaseConnConf,
		tableName:        ItemsTableName,
	}
}

type Item struct {
	Id   int64   `json:"id"`
	Name string  `json:"name"`
	Num  float64 `json:"num,omitempty"`
}

func (dbi *itemsDbClientImpl) SelectAllItems() ([]Item, error) {
	dbConn, err := dbOpenConnection(dbi.databaseConnConf)
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	sqlQuery := fmt.Sprintf("SELECT id, name FROM %s", dbi.tableName)
	result, err := dbConn.Query(sqlQuery)
	if err != nil {
		log.Printf("failed to execute query: %s", err.Error())
		return nil, err
	}
	defer result.Close()

	var items []Item
	for result.Next() {
		var item Item
		if err := result.Scan(&item.Id, &item.Name); err != nil {
			log.Printf("failed to read query result: %s", err.Error())
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (dbi *itemsDbClientImpl) GetItemById(id string) (*Item, error) {
	if !isIntNumber(id) {
		return nil, fmt.Errorf("%s is not valid id", id)
	}

	dbConn, err := dbOpenConnection(dbi.databaseConnConf)
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	sqlQuery := fmt.Sprintf("SELECT * FROM %s WHERE id = %s", dbi.tableName, id)
	result := dbConn.QueryRow(sqlQuery)

	item := &Item{}
	if err := result.Scan(&item.Id, &item.Name, &item.Num); err != nil {
		if err == sql.ErrNoRows {
			// There is no record with given ID
			return nil, nil
		}
		log.Printf("failed to execute query: %s", err.Error())
		return nil, err
	}
	return item, nil
}

func (dbi *itemsDbClientImpl) InsertItem(item Item) (string, error) {
	dbConn, err := dbOpenConnection(dbi.databaseConnConf)
	if err != nil {
		return "", err
	}
	defer dbConn.Close()

	sqlQuery := fmt.Sprintf("INSERT INTO %s(name, num) VALUES ('%s', %f)", dbi.tableName, item.Name, item.Num)
	result, err := dbConn.Exec(sqlQuery)
	if err != nil {
		log.Printf("failed to execute query: %s", err.Error())
		return "", err
	}
	newItemId, err := result.LastInsertId()
	if err != nil {
		log.Printf("failed to get inserted item ID: %s", err.Error())
		return "", err
	}
	return fmt.Sprint(newItemId), nil
}

func (dbi *itemsDbClientImpl) DeleteItemById(id string) (bool, error) {
	if !isIntNumber(id) {
		return false, fmt.Errorf("%s is not valid id", id)
	}

	dbConn, err := dbOpenConnection(dbi.databaseConnConf)
	if err != nil {
		return false, err
	}
	defer dbConn.Close()

	sqlQuery := fmt.Sprintf("DELETE FROM %s	WHERE id = %s", dbi.tableName, id)
	result, err := dbConn.Exec(sqlQuery)
	if err != nil {
		log.Printf("failed to execute query: %s", err.Error())
		return false, err
	}
	deletedRows, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get number of deleted rows: %s", err.Error())
		return false, err
	}

	return deletedRows > 0, nil
}

func (dbi *itemsDbClientImpl) GetRandomItem() (*Item, error) {
	dbConn, err := dbOpenConnection(dbi.databaseConnConf)
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	sqlQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY RAND() LIMIT 1", dbi.tableName)
	result := dbConn.QueryRow(sqlQuery)

	item := &Item{}
	if err := result.Scan(&item.Id, &item.Name, &item.Num); err != nil {
		if err == sql.ErrNoRows {
			// Empty table
			return nil, nil
		}
		log.Printf("failed to execute query: %s", err.Error())
		return nil, err
	}
	return item, nil
}

func isIntNumber(num string) bool {
	if _, err := strconv.Atoi(num); err == nil {
		return true
	}
	return false
}
