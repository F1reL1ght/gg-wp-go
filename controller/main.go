package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type RegisterationRequest struct {
	Device string `json:"device"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
}
type Color struct {
	R int
	G int
	B int
}

type Device struct {
	name  string
	ip    string
	port  int
	alive bool
	color Color
}

var devices map[string]*Device
type HttpResponse map[string]interface{}

func _register(w http.ResponseWriter, r *http.Request) {
	var request RegisterationRequest
	 body ,err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(body, &request)
	devices[request.IP] = &Device{name: request.Device, ip: request.IP, port: request.Port, alive: false, color: Color{} }
	w.Header().Set("Content-Type", "application/json")
	responseBody, err := json.Marshal(HttpResponse{"code": 200, "details": "registered successfully", "error": "" })
	w.Write(responseBody)
}

func _devices(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "%v", devices)
}

func _setRGB(w http.ResponseWriter, r *http.Request) {

}

func main() {
	go _health()
	http.HandleFunc("/register", _register)
	http.HandleFunc("/devices", _devices)
	http.HandleFunc("/set-rgb", _setRGB)



	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}
}


func _health() {
	for {
		time.Sleep(10 * time.Second)
		for _, device := range devices {
			url := fmt.Sprintf("http://%v:%v/health", device.ip, device.port )
			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			respBody, err := ioutil.ReadAll(resp.Body)
			fmt.Println((string)(respBody))
		}
	}
}

func init() {
	//TODO: Check if the desk, group, deskgroup tables are created if not create them
	db, err := sql.Open("sqlite3", "./ggwp.db")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.Exec(`PRAGMA foreign_keys=ON`)
	tblDesk, err := db.Prepare("CREATE TABLE IF NOT EXISTS desk (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, ip TEXT, port TEXT, red INTEGER, green INTEGER, blue INTEGER, health TEXT, retries INTEGER)")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	tblDesk.Exec()
	tblGroup, err := db.Prepare("CREATE TABLE IF NOT EXISTS collection (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	tblGroup.Exec()
	db.Exec("CREATE TABLE IF NOT EXISTS deskcollection(id INTEGER PRIMARY KEY AUTOINCREMENT, collection_id INTEGER, desk_id INTEGER FOREIGN KEY collection_id REFERENCES collection(id))")
	fmt.Println("func[init] initilization complete")
	}