package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PORT  = ":8080"
	INDEX = "templates/index.html"
)

func main() {

	// Enable Kubernetes cluster mode if running in a cluster (disabled by default)
	clusterModePtr := flag.Bool("cluster", false, "Enable kubernetes cluster mode")

	flag.Parse()

	if *clusterModePtr {
		log.Println("Enabling kubernetes cluster mode...")
		go watcher()
	}

	log.Printf("Started Ethereum Data Fetcher on port [%v] ", PORT)

	router := mux.NewRouter()

	router.HandleFunc("/", RootHandler)

	router.HandleFunc("/query", QueryHandler)

	router.HandleFunc("/api", ApiHandler).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
