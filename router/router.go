package router

import (
	"github.com/gorilla/mux"
	"github.com/i4ba1/CustomerOrderAPI/middleware"
)

// SetupRouter is exported and used in main.go
func SetupRouter() *mux.Router {
	newRouter := mux.NewRouter()
	newRouter.HandleFunc("/api/register", middleware.CreateUser).Methods("POST")
	newRouter.HandleFunc("/api/login", middleware.Login).Methods("POST")
	newRouter.HandleFunc("/api/refresh", middleware.Refresh).Methods("POST")
	newRouter.HandleFunc("/api/createOrder", middleware.CreateOrder).Methods("POST")
	return newRouter
}
