package main

import (
	"database/sql"
	"github.com/pkg/errors"
	"fmt"
)

var (
	DATABASE_NOT_OPEN = errors.New("database not open")
)

type sqlite struct {
	db *sql.DB
}

func (driver *sqlite) Open(database string) (error) {
	db, err := sql.Open("sqlite3", database)
	driver.db = db
	return err
}

func (driver *sqlite) EnableForeignKeys() (error) {
	if driver.db == nil {
		return DATABASE_NOT_OPEN
	}
	driver.db.Exec(`PRAGMA foreign_keys=ON`)
	return nil
}

func (driver sqlite) Query(params ...interface{}) (interface{}, error) {
	if driver.db == nil {
		return nil, DATABASE_NOT_OPEN
	}
	queryStr := params[0].(string)
	result, err := driver.db.Exec(queryStr)
	return result, err
}

//TODO: Guard against SQL Injection
func (driver sqlite) Insert(params ...interface{}) (interface{}, error) {
	queryStr := params[0].(string)
	stmt, err := driver.db.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	return stmt.Exec()
}

func (driver sqlite) Delete(params ...interface{}) (interface{}, error) {
	queryStr := params[0].(string)
	stmt, err := driver.db.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	return stmt.Exec()
}

func (driver sqlite) AddDesk(desk Desk) (bool, error) {
	sqlQuery := `INSERT INTO ` + TABLE_DESK + `(name, ip, port, red, green, blue, health, retries) VALUES (` +
	    desk.Name + `,` +
		desk.IP + `,` +
		desk.Port + `,` +
	    fmt.Sprintf("%v", desk.Red) +
		fmt.Sprintf("%v", desk.Green) +
		fmt.Sprintf("%v", desk.Blue) +
		desk.Health +
		fmt.Sprintf("%v", desk.Retries) + `,`+`)`
		_, err := driver.Insert(sqlQuery)
		if err != nil {
			return false, err
		}
		return true, nil
}

func (driver sqlite) AddCollection(collection Collection) (bool, error) {
	sqlQuery := `INSERT INTO ` + TABLE_COLLECTION + `(name) VALUES (` + collection.Name + `)`
	_, err := driver.Insert(sqlQuery)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (driver sqlite) AddDeskToCollection(deskID int, collectionID int) (bool, error) {
	sqlQuery := `INSERT INTO ` + TABLE_COLLECTION + `(collection_id, desk_id) VALUES (` + fmt.Sprintf("%v", deskID) + `,` + fmt.Sprintf("%v", collectionID) + `)`
	_, err := driver.Insert(sqlQuery)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* Delete the Association from the Database*/
func (driver sqlite) RemoveDeskFromCollection(associationID int) (bool, error) {
	sqlQuery := `DELETE FROM ` + TABLE_DESK_COLLECTION + ` WHERE id = ` + fmt.Sprintf("%v", associationID)
	_, err := driver.Delete(sqlQuery)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (driver sqlite) GetDeskCollections() (*DeskCollection, error) {
	sqlQuery := `SELECT * FROM ` + TABLE_COLLECTION + ` INNER JOIN ` + TABLE_DESK_COLLECTION + ` ON ` + TABLE_COLLECTION + `.id = ` + TABLE_DESK_COLLECTION + `.collection_id` +
				` INNER JOIN ` + TABLE_DESK + ` ON ` + TABLE_DESK + `.id = ` + TABLE_DESK_COLLECTION + `.desk_id`

	fmt.Println(sqlQuery)
	return nil, nil
}