package main

import "fmt"
import "net/http"
import "github.com/rakshitg600/notakto-solo/routes"


func main() {
	routes.RegisterRoutes()
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
