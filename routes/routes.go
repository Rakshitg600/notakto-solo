package routes

import (
	"net/http"
	"github.com/rakshitg600/notakto-solo/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/", handlers.HelloHandler)
	http.HandleFunc("/create",handlers.CreateHandler)
}
