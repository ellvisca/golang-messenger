package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ellvisca/messenger/controllers"
	"github.com/gorilla/mux"
	"github.com/maple-ai/syrup"
)

func main() {
	router := syrup.New(mux.NewRouter())

	// User router
	router.Post("/api/v1/user", controllers.CreateClient)

	// Login router
	router.Post("/api/v1/user/login", controllers.ClientLogin)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Listening on port ", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
