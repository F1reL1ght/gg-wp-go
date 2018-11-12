package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	//"fmt"
	"fmt"
)

type Desk struct {
	Name string `json:"name"`
	IP string	`json:"ip"`
	Port string `json:"port"`
	Red int 	`json:"red"`
	Green int	`json:"green"`
	Blue int	`json:"blue"`
	Health string `json:"health"`
	Retries int	`json:"retries"`
}

type Collection struct {
	Name string
}

type DeskCollection struct {
	Name string
	Desks []Desk
}

const (
	TABLE_DESK = "desk"
	TABLE_COLLECTION = "collection"
	TABLE_DESK_COLLECTION = "deskcollection"
)

/*Add Desk to Persistent Monitoring*/
func _addDesk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var responseBody []byte
	body ,err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var desk Desk
	err = json.Unmarshal(body, &desk)

	if err != nil {
		responseBody, _ = json.Marshal(HttpResponse{"code": 500, "details": http.StatusInternalServerError, "error": err.Error() })
		w.Write(responseBody)
		return
	}
	desk.Health = "false"
	_, err = addDesk(desk)

	if err == nil {
		responseBody, _ = json.Marshal(HttpResponse{"code": 200, "details": "added Desk successfully", "error": "" })
	} else {
		responseBody, _ = json.Marshal(HttpResponse{"code": 500, "details": http.StatusInternalServerError, "error": err.Error() })
	}
	w.Write(responseBody)
}

/*Add Collection*/
func _addCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var responseBody []byte
	body ,err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var collection Collection
	err = json.Unmarshal(body, &collection)
	if err != nil {
		responseBody, _ = json.Marshal(HttpResponse{"code": 500, "details": http.StatusInternalServerError, "error": err.Error() })
		w.Write(responseBody)
		return
	}
	_, err = addCollection(collection)
	if err == nil {
		responseBody, _ = json.Marshal(HttpResponse{"code": 200, "details": "added Collection successfully", "error": "" })
	} else {
		fmt.Println(err)
		responseBody, _ = json.Marshal(HttpResponse{"code": 500, "details": http.StatusInternalServerError, "error": err.Error() })
	}
	w.Write(responseBody)
}

func _getCollectionList(w http.ResponseWriter, r *http.Request) {
	getCollections()
}

func addDesk(desk Desk) (bool, error) {
	idx, err := (database).AddDesk(desk)
	if err != nil {
		return false, err
	}
	return  (database).AddDeskToCollection(idx, 1)
}

func addCollection(collection Collection) (bool, error) {
	return (database).AddCollection(collection)
}

/* Add Desk to Collection */
func addDeskToCollection(deskID int64, collectionID int64) (bool, error) {
	return (database).AddDeskToCollection(deskID, collectionID)
}

/* gets all Collections and Associated Desks */
func getCollections() () {
	fmt.Println((database).GetDeskCollections())
}