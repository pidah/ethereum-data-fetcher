package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ethereum struct {
	EthereumAddress string `json:"ethereum_address"`
}

func RootHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("RootHandler started and redirecting to the index page")

	t, _ := template.ParseFiles(INDEX)

	t.Execute(w, nil)

}

func QueryHandler(w http.ResponseWriter, r *http.Request) {

	//Retrieve the HTML form parameter of POST method
	e := r.FormValue("ethereum-data")

	t, err := template.ParseFiles(INDEX)

	if err != nil {
		fmt.Println(err.Error())
	}

	g := GetEthereumData(e)
	fmt.Println(g)

	t.Execute(w, g)

}

func ApiHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var request ethereum

	err := decoder.Decode(&request)

	g := GetEthereumData(request.EthereumAddress)

	jsonData, err := json.Marshal(g)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprint(w, string(jsonData))
	return
}
