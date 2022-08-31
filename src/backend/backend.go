package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB     *sql.DB
	Port   string
	Router *mux.Router
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func (a *App) Initialize() {
	DB, err := sql.Open("sqlite3", "../../practiceit.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	a.DB = DB
	a.Router = mux.NewRouter()
	a.initalizeRoutes()
}

func (a *App) initalizeRoutes() {
	a.Router.HandleFunc("/products", a.allProducts).Methods("GET")
	a.Router.HandleFunc("/product/{id}", a.fetchProduct).Methods("GET")
	a.Router.HandleFunc("/products", a.newProduct).Methods("POST")

	a.Router.HandleFunc("/orders", a.allOrders).Methods("GET")
	a.Router.HandleFunc("/order/{id}", a.fetchOrder).Methods("GET")
	a.Router.HandleFunc("/orders", a.newOrder).Methods("POST")
	a.Router.HandleFunc("/orderitems", a.newOrderItems).Methods("POST")
}

func (a *App) allProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(a.DB)
	if err != nil {
		fmt.Printf("getProducts error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) fetchProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p product
	p.ID, _ = strconv.Atoi(id)
	err := p.getProduct(a.DB)
	if err != nil {
		fmt.Printf("getProduct error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) newProduct(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var p product
	json.Unmarshal(reqBody, &p)

	err := p.createProduct(a.DB)
	if err != nil {
		fmt.Printf("newProduct error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) allOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := getOrders(a.DB)
	if err != nil {
		fmt.Printf("allOrders error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, orders)
}

func (a *App) fetchOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var o order
	o.ID, _ = strconv.Atoi(id)
	err := o.getOrder(a.DB)
	if err != nil {
		fmt.Printf("getOrder error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, o)
}

func (a *App) newOrder(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var o order
	json.Unmarshal(reqBody, &o)

	err := o.createOrder(a.DB)
	if err != nil {
		fmt.Printf("newOrder error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, item := range o.Items {
		var oi orderItem
		oi = item
		oi.OrderID = o.ID
		err := oi.createOrderItem(a.DB)
		if err != nil {
			fmt.Printf("newOrder, newOrderItem error: %s\n", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, o)
	}
}

func (a *App) newOrderItems(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var ois []orderItem
	json.Unmarshal(reqBody, &ois)

	for _, item := range ois {
		var oi orderItem
		oi = item
		err := oi.createOrderItem(a.DB)
		if err != nil {
			fmt.Printf("newOrderItem error: %s\n", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, ois)
}

func (a *App) Run() {
	fmt.Println("Server started and listening on port ", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

//Helper functions
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
