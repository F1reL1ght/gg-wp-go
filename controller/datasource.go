package main

type datasource interface {
	//Query(params ...interface{} ) (interface{}, error)
	AddDesk(desk Desk) (bool, error)
	AddCollection(collection Collection) (bool, error)
	/*Removes Desk from Group*/
	//RemoveDesk(deskID int, collectionID int) (bool, error)
	/*Removes Group and All associations to the group*/
	//RemoveCollection(collectionID int) (bool, error)
	AddDeskToCollection(deskID int, collectionID int) (bool, error)
	RemoveDeskFromCollection(associationID int) (bool, error)
}