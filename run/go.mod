module example.com/run

go 1.19

replace example.com/backend => ../router

require example.com/backend v0.0.0-00010101000000-000000000000

require github.com/gorilla/mux v1.8.0 // indirect
