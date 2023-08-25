package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codemunsta/risevest-test/src/routers"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Server Starting Up")

	router := routers.NewRouter()
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	c.Handler(router)
}
