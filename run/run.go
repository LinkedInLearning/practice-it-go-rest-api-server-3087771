package main

import (
	"example.com/backend"
)

func main() {
	// check if its because of capital letters
	myrouterstructobj := backend.App{}
	myrouterstructobj.Port = ":8080"
	myrouterstructobj.Initialize()
	myrouterstructobj.Run()
}
