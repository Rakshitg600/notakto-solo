package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	db "github.com/rakshitg600/notakto-solo/db/generated"
	"github.com/rakshitg600/notakto-solo/handlers"
	"github.com/rakshitg600/notakto-solo/routes"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	queries := db.New(conn)
	handler := handlers.NewHandler(queries)

	e := echo.New()
	// ✅ Enable CORS for frontend (Next.js at localhost:3000)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://deploy-preview-355--staging-notakto.netlify.app/", "https://deploy-preview-356--staging-notakto.netlify.app/"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	routes.RegisterRoutes(e, handler)
	e.Logger.Fatal(e.Start(":1323"))
}
