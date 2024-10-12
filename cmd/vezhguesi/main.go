package main

import (
	"fmt"
	"net/http"

	_ "vezhguesi/cmd/vezhguesi/docs" // Import the generated docs package
	db "vezhguesi/core/db"
	usersvc "vezhguesi/core/users"

	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	db.AutoMigrate(&usersvc.User{})

	// Serve Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Start the server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
