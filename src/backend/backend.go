package backend

import (
	"fmt"
	"log"
	"net/http"
)

func dimension(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Lenght X Wide X Height")
}

func Run(addr string) {
	http.HandleFunc("/", dimension)
	fmt.Println("Server Has Started On Port ", addr)
	log.Fatal(http.ListenAndServer(addr, nil))

}
