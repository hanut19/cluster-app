package main

import (
	"cluster-app/db"
	"cluster-app/handlers"
	"cluster-app/middleware"
	"log"
	"net/http"
)

func main() {
	err := db.InitDB("config.json")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handlers.Login)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/portal", middleware.RequireLogin(handlers.Portal))
	http.HandleFunc("/update", middleware.RequireLogin(middleware.RequireAdmin(handlers.Update)))

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
