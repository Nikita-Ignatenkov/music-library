package controllers

import (
	"encoding/json"
	"io"
	"log"
	"music-library/app/database"
	"music-library/app/models"
	"net/http"
	"net/url"
)

func GetSongs(w http.ResponseWriter, r *http.Request) {
	var songs []models.Song
	if err := database.DB.Find(&songs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}
func AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	groupEscaped := url.QueryEscape(song.Group)
	songEscaped := url.QueryEscape(song.Name)
	apiURL := "http://localhost:8081/info?group=" + groupEscaped + "&song=" + songEscaped
	log.Printf("Making request to mock API at %s\n", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Failed to call external API: %s\n", err)
		http.Error(w, "Failed to call external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Non-200 response received: %d, body: %s\n", resp.StatusCode, string(bodyBytes))
		http.Error(w, "API request returned non-200 status", http.StatusInternalServerError)
		return
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		return
	}

	song.ReleaseDate = songDetail.ReleaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	if err := database.DB.Create(&song).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}
