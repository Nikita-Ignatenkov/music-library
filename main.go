package main

import (
	"encoding/json"
	"log"
	"music-library/app/config"
	"music-library/app/database"
	"music-library/app/routes"
	"net/http"
	"sync"
)

func main() {
	log.Println("INFO: Starting the music library application...")

	config.LoadConfig()
	database.ConnectDatabase()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		router := routes.RegisterRoutes()

		log.Println("INFO: Server started at :8000")
		if err := http.ListenAndServe(":8000", router); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		startMockAPIServer()
	}()

	wg.Wait()
	log.Println("INFO: Shutting down the application.")
}

// мой тестовый Mock API запускается на порту 8081
func startMockAPIServer() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		song := r.URL.Query().Get("song")

		log.Printf("DEBUG: Received request for group: %s, song: %s\n", group, song)

		if group == "" || song == "" {
			log.Println("Group or song is empty, returning 400 Bad Request")
			http.Error(w, "Group or song is required", http.StatusBadRequest)
			return
		}

		var response struct {
			ReleaseDate string `json:"releaseDate"`
			Text        string `json:"text"`
			Link        string `json:"link"`
		}

		if group == "Muse" && song == "Supermassive Black Hole" {
			response = struct {
				ReleaseDate string `json:"releaseDate"`
				Text        string `json:"text"`
				Link        string `json:"link"`
			}{
				ReleaseDate: "2006-07-16",
				Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
				Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
			}
		} else if group == "ласковый май" && song == "белые розы" {
			response = struct {
				ReleaseDate string `json:"releaseDate"`
				Text        string `json:"text"`
				Link        string `json:"link"`
			}{
				ReleaseDate: "1988-02-16",
				Text:        "Белые pозы, белые pозы, беззащитны шипы",
				Link:        "https://youtu.be/CTpyz63q-6c?si=3GfsZwpV6EU8qJTk",
			}
		} else {
			log.Printf("Song not found for group: %s, song: %s\n", group, song)
			http.Error(w, "Song not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Mock API server started at :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Failed to start mock API server:", err)
	}
}
