package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codemunsta/risevest-test/src/db"
	"github.com/codemunsta/risevest-test/src/routers"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var port = envPortOr("3000")

	db.InitDB()

	fmt.Println("Server Starting Up")

	router := routers.NewRouter()
	http.Handle("/", router)
	http.ListenAndServe(port, nil)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	c.Handler(router)
}

func envPortOr(port string) string {
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	return ":" + port
}
