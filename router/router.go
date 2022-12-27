package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Myrouterobj *mux.Router
	Port        string
}

func (a *App) Initialize() {

	a.Myrouterobj = mux.NewRouter()

}

func (a *App) Run() {
	a.Myrouterobj.HandleFunc("/", getRequest).Methods("GET")
	a.Myrouterobj.HandleFunc("/", postRequest).Methods("POST")
	http.Handle("/", a.Myrouterobj)
	log.Fatal(http.ListenAndServe(a.Port, nil))
	fmt.Println("Server started on the port ", a.Port)

}

func getRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a GET")
}

func postRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a POST")
}
