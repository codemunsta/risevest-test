package routers

import (
	"net/http"

	h "github.com/codemunsta/risevest-test/src/handlers"
	mWare "github.com/codemunsta/risevest-test/src/middleware"
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

	// files and folder
	router.HandleFunc("/api/user/upload/file", mWare.IsAuthenticated(h.UploadFile)).Methods(http.MethodPost)
	router.HandleFunc("/api/user/create/folder", mWare.IsAuthenticated(h.CreateFolder)).Methods(http.MethodPost)
	router.HandleFunc("/api/user/download", h.FileDownload).Methods(http.MethodGet)
	router.HandleFunc("/api/user/get/folders", mWare.IsAuthenticated(h.GetFolders)).Methods(http.MethodGet)
	router.HandleFunc("/api/user/get/files", mWare.IsAuthenticated(h.GetFiles)).Methods(http.MethodGet)
	return router
}
