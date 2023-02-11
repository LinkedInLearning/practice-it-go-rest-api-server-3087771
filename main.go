package main

import (
	"fmt"
	"log"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func main() {
	http.HandleFunc("/", HelloWorld)
	fmt.Println("Server Started and Listening On Localhost: 9003")
	log.Fatal(http.ListenAndServe(":9003", nil))
}
