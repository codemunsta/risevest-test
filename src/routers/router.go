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
	router.HandleFunc("/api/get_user", h.TestUserAuth).Methods(http.MethodGet)

	// registeration
	router.HandleFunc("/api/user/register", h.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/admin/register", h.RegisterAdmin).Methods(http.MethodPost)

	// login
	router.HandleFunc("/api/user/login", h.LoginAuthentication).Methods(http.MethodPost)
	router.HandleFunc("/api/admin/login", h.LoginAdmin).Methods(http.MethodPost)
	return router
}
