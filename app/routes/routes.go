package routes

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"music-library/app/controllers"
	_ "music-library/docs"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/songs", controllers.GetSongs).Methods("GET")
	router.HandleFunc("/songs/{id}/text", controllers.GetSongTextWithPagination).Methods("GET")
	router.HandleFunc("/songs", controllers.AddSong).Methods("POST")
	router.HandleFunc("/songs/{id}", controllers.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", controllers.DeleteSong).Methods("DELETE")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}
