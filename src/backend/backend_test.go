package backend_test

import (
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"example.com/backend"
	"log"
	"net/http"
	"net/http/httptest"
)

var a backend.App

const tableProductCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
	id INT NOT NULL PRIMARY KEY AUTOINCREMENT,
	productCode VARCHAR(25) NOT NULL,
	name VARCHAR(256) NOT NULL,
	inventory INT NOT NULL,
	price INT NOT NULL,
	status VARCHAR(64) NOT NULL
)`

func TestMain(m *testing.M) {
	a = backend.App{}
	a.Initialize()
	ensureTableExists()
	code := m.Run()

	clearProductTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableProductCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearProductTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'products'")
}

func TestGetNonExistentProduct(t *testing.T) {
	clearProductTable()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "sql: no rows in result set" {
		t.Errorf("Expected the 'error' key of the response to be set to 'sql:no rows in result set'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearProductTable()

	payload := []byte(`{"productCode":"TEST12345", "name":"ProductTest", "inventory":1, "price":1, "status":"testing"}`)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["productCode"] != "TEST12345" {
		t.Errorf("Expected productCode to be 'TEST12345'. Got '%v'", m["productCode"])
	}

	if m["name"] != "ProductTest" {
		t.Errorf("Expected name to be 'ProductTest'. Got '%v'", m["name"])
	}

	if m["inventory"] != 1.0 {
		t.Errorf("Expected inventory to be '1'. Got '%v'", m["inventory"])
	}

	if m["price"] != 1.0 {
		t.Errorf("Expected price to be '1'. Got '%v'", m["price"])
	}

	if m["status"] != "testing" {
		t.Errorf("Expected status to be 'testing'. Got '%v'", m["status"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected id to be '1'. Got '%v'", m["id"])
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
