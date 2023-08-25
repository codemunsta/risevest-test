package routers

import (
	"net/http"

	h "github.com/codemunsta/risevest-test/src/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	router.HandleFunc("/api/ping", h.NewApi).Methods(http.MethodGet)
	return router
}
