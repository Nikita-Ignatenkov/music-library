package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"music-library/app/database"
	"music-library/app/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: Received request to get songs")
	var songs []models.Song
	group := r.URL.Query().Get("group")
	name := r.URL.Query().Get("song")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var limit, offset int
	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Println("INFO: Invalid limit value")
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			log.Println("INFO: Invalid offset value")
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
	}

	query := database.DB.Model(&songs)

	if group != "" {
		query = query.Where("artist = ?", group)
	}
	if name != "" {
		query = query.Where("name = ?", name)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&songs).Error; err != nil {
		log.Println("INFO: Failed to retrieve songs")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func GetSongTextWithPagination(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Println("DEBUG: Received request for song text with ID:", id)
	var song models.Song

	if err := database.DB.First(&song, id).Error; err != nil {
		log.Println("INFO: Song not found with ID:", id)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	verses := strings.Split(song.Text, "\n\n")
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	perPage := 2

	start := (page - 1) * perPage
	end := start + perPage

	if start >= len(verses) {

		log.Println("INFO: Page exceeds total number of verses")
		http.Error(w, "Page exceeds total number of verses", http.StatusBadRequest)
		return
	}

	if end > len(verses) {
		end = len(verses)
	}

	response := verses[start:end]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Println("DEBUG: Received request to delete song with ID:", id)
	if err := database.DB.Delete(&models.Song{}, id).Error; err != nil {
		log.Println("INFO: Failed to delete song with ID:", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("DEBUG: Successfully deleted song with ID:", id)
	w.WriteHeader(http.StatusNoContent)
}

func UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Println("INFO: Failed to decode request body for update")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Model(&models.Song{}).Where("id = ?", id).Updates(song).Error; err != nil {
		log.Println("INFO: Failed to update song with ID:", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("DEBUG: Successfully updated song with ID:", id)
	w.WriteHeader(http.StatusNoContent)
}

func AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Println("INFO: Failed to decode request body for adding song")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	apiURL := "http://localhost:8081/info?group=" + url.QueryEscape(song.Group) + "&song=" + url.QueryEscape(song.Name)
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("INFO: Failed to call external API or received non-200 status")
		http.Error(w, "Failed to call external API", http.StatusInternalServerError)
		return
	}

	if err := database.DB.Create(&song).Error; err != nil {
		log.Println("INFO: Failed to save song to the database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("DEBUG: Successfully added song:", song.Name)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}
