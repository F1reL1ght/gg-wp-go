package main

import (
	"database/sql"
	"github.com/pkg/errors"
	"fmt"
)

var (
	DATABASE_NOT_OPEN = errors.New("database not open")
	DATABASE_ERROR = errors.New("cannot open database")
)

type sqlite struct {
	db *sql.DB
}

func(driver *sqlite) init() error{
	driver.db.Exec(`PRAGMA foreign_keys=ON`)
	_, err := driver.db.Exec("CREATE TABLE IF NOT EXISTS desk (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, ip TEXT, port TEXT, red INTEGER, green INTEGER, blue INTEGER, health TEXT, retries INTEGER)")
	if err != nil {
		fmt.Println(err.Error())
		return  DATABASE_ERROR
	}
	_, err = driver.db.Exec("CREATE TABLE  collection (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)")
	if err == nil {
		if _, err := driver.AddCollection(Collection{Name: "All"}); err != nil {
			return DATABASE_ERROR
		}
	}

	_, err = driver.db.Exec("CREATE TABLE IF NOT EXISTS deskcollection(id INTEGER PRIMARY KEY AUTOINCREMENT, collection_id INTEGER, desk_id INTEGER, CONSTRAINT fk_collection FOREIGN KEY (collection_id) REFERENCES collection(id), CONSTRAINT fk_desk FOREIGN KEY (desk_id) REFERENCES desk(id))")
	if err != nil {
		fmt.Println(err.Error())
		return  DATABASE_ERROR
	}

	return nil
}

func NewSQLiteDriver(database string) (*sqlite, error) {
	driver := sqlite{}
	if err := driver.open(database); err != nil {
		return nil, DATABASE_ERROR
	}
	return &driver, nil
}



func (driver *sqlite) open(database string) (error) {
	db, err := sql.Open("sqlite3", database)
	driver.db = db
	return err
}

//func (driver *sqlite) EnableForeignKeys() (error) {
//	if driver.db == nil {
//		return DATABASE_NOT_OPEN
//	}
//
//	return nil
//}

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
	return driver.db.Exec(queryStr)

}

func (driver sqlite) Delete(params ...interface{}) (interface{}, error) {
	queryStr := params[0].(string)
	stmt, err := driver.db.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	return stmt.Exec()
}

func (driver *sqlite) AddDesk(desk Desk) (int64, error) {
	sqlQuery := `INSERT INTO ` + TABLE_DESK + `(name, ip, port, red, green, blue, health, retries) VALUES (` +
	    `'`+ desk.Name +`'` + `,` +
		`'` + desk.IP + `'` +  `,` +
		`'` + desk.Port +`'` + `,` +
	    fmt.Sprintf("%v", desk.Red) + `,` +
		fmt.Sprintf("%v", desk.Green) + `,` +
		fmt.Sprintf("%v", desk.Blue) + `,` +
		`'` + desk.Health + `'` + `,` +
		fmt.Sprintf("%v", desk.Retries) +`)`
		fmt.Println(sqlQuery)
		result, err := driver.Insert(sqlQuery)
		if err != nil {
			return 0, err
		}
		r := result.(sql.Result)
		idx, err := r.LastInsertId()
		return idx, err
}

func (driver *sqlite) AddCollection(collection Collection) (bool, error) {
	sqlQuery := `INSERT INTO ` + TABLE_COLLECTION + `(name) VALUES ('` + collection.Name + `')`
	//fmt.Println(sqlQuery)
	result, err := driver.Insert(sqlQuery)
	r, err := result.(sql.Result).RowsAffected()
	fmt.Printf("%v",r )
	if err != nil {
		return false, err
	}
	return true, nil
}



func (driver sqlite) AddDeskToCollection(deskID int64, collectionID int64) (bool, error) {
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