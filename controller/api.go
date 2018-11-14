package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	//"fmt"

	"database/sql"
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

type SQLDesk struct {
	DeskID string `json:"desk_id"`
	Name string `json:"desk_name"`
	IP string	`json:"desk_ip"`
	Port string `json:"desk_port"`
	Red string 	`json:"red"`
	Green string	`json:"green"`
	Blue string	`json:"blue"`
	Health string `json:"desk_health"`
	Retries string	`json:"retries"`
}

type SQLCollection struct {
	CollectionID string `json:"collection_id"`
	Name string `json:"collection_name"`
}

type Collection struct {
	Name string
}

type SQLDeskCollection struct {
	SQLCollection
	Desks []SQLDesk `json:"desks"`
}

type AssociateRequest struct {
	deskID int	`json:"desk_id"`
	collectionID int  `json:"collection_id"`
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
		responseBody, _ = json.Marshal(HttpResponse{"code": 500, "details": http.StatusInternalServerError, "error": err.Error() })
	}
	w.Write(responseBody)
}

func _getCollectionList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var responseBody []byte
	if collection, err := getCollections(); err != nil {
		responseBody, _ = json.Marshal(HttpResponse{"code": 500, "details": http.StatusInternalServerError, "error": err.Error() })
	} else {
		responseBody, err = json.Marshal(HttpResponse{"code": 200, "collections": collection, "error": "" })
		if err != nil {
			panic(err)
		}
	}

	w.Write(responseBody)
}

func _addDeskToCollection(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//var responseBody []byte
	//body ,err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	panic(err)
	//}
	//var request AssociateRequest
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
func getCollections() ([]*SQLDeskCollection, error) {

	builder := make(map[string]*SQLDeskCollection)
	result, err := (database).GetDeskCollections()

	if err != nil {
		panic(err)
	}
	r := result.(*sql.Rows)
	defer r.Close()

	var DeskID, Name, IP, Port, Red, Green, Blue, Health, Retries, CollectionID, AssociationID, Associated_Desk_ID, Associated_Collection_ID, CollectionName string

	for r.Next() {

		err := r.Scan(&CollectionID,
					  &CollectionName,
					  &AssociationID,
					  &Associated_Collection_ID,
					  &Associated_Desk_ID,
					  &DeskID,
					  &Name,
					  &IP,
					  &Port,
					  &Red,
					  &Green,
					  &Blue,
					  &Health,
					  &Retries)

		if err != nil {
			panic(err)
		}

		if dc := builder[CollectionName]; dc == nil {
			var deskCollection SQLDeskCollection
			deskCollection.CollectionID = CollectionID
			deskCollection.Name = CollectionName
			deskCollection.Desks = append(deskCollection.Desks, SQLDesk{DeskID, Name, IP, Port, Red, Green, Blue, Health, Retries})
			builder[CollectionName] = &deskCollection
			} else {
			dc.Desks = append(dc.Desks, SQLDesk{DeskID, Name, IP, Port, Red, Green, Blue, Health, Retries})
		}
	}

	var collections []*SQLDeskCollection
	for _, collection := range builder {
			collections = append(collections, collection)
	}

	return collections, nil
}