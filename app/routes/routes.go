package routes

import (
	"github.com/gorilla/mux"
	"music-library/app/controllers"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/songs", controllers.GetSongs).Methods("GET")
	router.HandleFunc("/songs/{id}/text", controllers.GetSongTextWithPagination).Methods("GET")
	router.HandleFunc("/songs", controllers.AddSong).Methods("POST")
	router.HandleFunc("/songs/{id}", controllers.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", controllers.DeleteSong).Methods("DELETE")
	return router
}
