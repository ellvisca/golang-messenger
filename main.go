package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ellvisca/messenger/app"
	"github.com/ellvisca/messenger/controllers"
	"github.com/gorilla/mux"
	"github.com/maple-ai/syrup"
)

func main() {
	router := syrup.New(mux.NewRouter())

	// Client router
	router.Post("/api/v1/client", controllers.CreateClient)

	// Login router
	router.Post("/api/v1/client/login", controllers.ClientLogin)

	// Hub router
	router.Post("/api/v1/hub", controllers.CreateHub)
	router.Post("/api/v1/hub/message", controllers.ReceiveMsg)

	// Middleware
	router.Use(app.JwtAuthentication)

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
