package router

import (
	"github.com/gorilla/mux"
	"github.com/i4ba1/CustomerOrderAPI/middleware"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/createUser", middleware.CreateUser).Methods("POST", "OPTIONS")
	return router
}
