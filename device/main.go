package main

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"time"
	//"os"
	"net"
	"log"
)
const MAX_RETRY = 10
var RETRY = 0
func main() {


	// isRegistered := register()
	//for !isRegistered {
	//	RETRY++
	//	time.Sleep(1000)
	//	if RETRY >= MAX_RETRY {
	//		os.Exit(-1)
	//	}
	//	isRegistered = register()
	//	fmt.Println(isRegistered)
	//}
	//addr, err :=  net.ResolveUDPAddr("udp:", ":8032")
	//fmt.Println("Here1")
	c, err := net.ListenPacket("udp", ":8032")
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()
	b := make([]byte, 512)

	n, peer, err := c.ReadFrom(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println((string)(b))
	fmt.Println(n, "bytes read from", peer)




	http.HandleFunc("/health", _health)
	if err := http.ListenAndServe(":8100", nil); err != nil {
		panic(err)
	}

}

func register()  bool {

	var registration = make(map[string]interface{},0)
	registration["device"] = "desk-1"
	registration["ip"] = "127.0.0.1"
	registration["port"] = 8100
	body, err := json.Marshal(registration)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:8090/register","application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return false
	}

	var HTTPResponse map[string]interface{}
	respBody, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody,&HTTPResponse)
	if err != nil {
		fmt.Printf("%v", err)
		return false
	}

	fmt.Println(HTTPResponse)
	//fmt.Println(reflect.TypeOf(HTTPResponse["code"]))

	if (HTTPResponse["code"]).(float64) != 200 {
		return false
	}
	return true
}

func _health(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"code": 200, "details": "healthy", "error": "" }
	responseBody, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	fmt.Println("health")
	w.Write(responseBody)
}

func _getRGB(w http.ResponseWriter, r *http.Request) {

}

func getRGB() {
	//TODO: Interface with the Hardware to determine Color
}

func _setRGB(w http.ResponseWriter, r *http.Request) {

}

func setRGB() {
	//TODO: Interface with the Hardware to Set Color
}