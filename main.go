package main

import (
	"fmt"
	"github.com/i4ba1/CustomerOrderAPI/router"
	"log"
	"net/http"
)

func main() {
	r := router.SetupRouter()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	fmt.Println("Starting server on the port 8787...")
	log.Fatal(http.ListenAndServe(":8787", r))
}