package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"Go-hashtable-Butterbrot/core" 
)



var rdb *redis.Client


func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", core.ServeHome).Methods("GET")
	r.HandleFunc("/api/put", core.PutHandler).Methods("POST")
	r.HandleFunc("/api/get/{key}", core.GetHandler).Methods("GET")
	r.HandleFunc("/api/delete/{key}", core.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/api/elements", core.GetAllElementsHandler).Methods("GET")

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

