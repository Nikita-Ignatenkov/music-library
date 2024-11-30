package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"music-library/app/database"
	"music-library/app/models"
)

// GetSongs возвращает список песен в зависимости от переданных параметров.
// @Summary Get a list of songs
// @Description Retrieve a list of songs based on optional group and song name filters, with pagination support using limit and offset parameters.
// @Accept json
// @Produce json
// @Param group query string false "Group name (artist)"
// @Param song query string false "Song name"
// @Param limit query int false "Limit the number of songs returned"
// @Param offset query int false "Offset the returned songs by this amount"
// @Success 200 {array} models.Song
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /songs [get]
func GetSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: Received request to get songs")
	var songs []models.Song

	group := r.URL.Query().Get("group")
	name := r.URL.Query().Get("song")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0
	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Println("INFO: Invalid limit value")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid limit",
			})
			return
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			log.Println("INFO: Invalid offset value")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid offset",
			})
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve songs",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// GetSongTextWithPagination возвращает текст песни с пагинацией по куплетам.
// @Summary Получение текста песни с пагинацией по куплетам
// @Description Получает текст песни по ее идентификатору и поддерживает пагинацию для возвращения определённого количества куплетов на странице.
// @Accept json
// @Produce json
// @Param id path string true "ID песни"
// @Param page query int false "Номер страницы для получения (индексация с единицы)"
// @Success 200 {array} string "Список куплетов"
// @Failure 400 {object} models.ErrorResponse "Неверный запрос, ошибка в параметрах"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Router /songs/{id}/text [get]
func GetSongTextWithPagination(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("DEBUG: Received request for song ID:", id)

	var song models.Song

	if err := database.DB.First(&song, id).Error; err != nil {
		log.Println("INFO: Song not found with ID:", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Song not found",
		})
		return
	}

	text := strings.ReplaceAll(song.Text, "\\n", "\n")
	verses := strings.Split(text, "\n\n")

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		log.Println("INFO: Invalid page value:", pageStr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Page must be a valid integer greater than 0",
		})
		return
	}

	perPage := 1000
	start := (page - 1) * perPage
	end := start + perPage

	if start >= len(verses) {
		log.Println("INFO: Page exceeds total number of verses")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Page exceeds total number of verses",
		})
		return
	}

	if end > len(verses) {
		end = len(verses)
	}

	response := verses[start:end]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteSong удаляет песню по идентификатору из базы данных.
// @Summary Удаление песни по ID
// @Description Удаляет песню из базы данных по указанному идентификатору.
// @Accept json
// @Produce json
// @Param id path string true "ID песни"
// @Success 204 "Песня успешно удалена"
// @Failure 400 {object} models.ErrorResponse "Ошибка в запросе"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id} [delete]
func DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("DEBUG: Received request to delete song with ID:", id)

	var song models.Song

	if err := database.DB.First(&song, id).Error; err != nil {
		log.Println("INFO: Song not found with ID:", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Song not found",
		})
		return
	}

	if err := database.DB.Delete(&song).Error; err != nil {
		log.Println("INFO: Failed to delete song with ID:", id)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete song",
		})
		return
	}

	log.Println("DEBUG: Successfully deleted song with ID:", id)
	w.WriteHeader(http.StatusNoContent)
}

// UpdateSong обновляет информацию о песне по идентификатору.
// @Summary Обновление песни по ID
// @Description Обновляет данные песни в базе данных по указанному идентификатору.
// @Accept json
// @Produce json
// @Param id path string true "ID песни"
// @Param song body models.Song true "Данные песни для обновления"
// @Success 204 "Песня успешно обновлена"
// @Failure 400 {object} models.ErrorResponse "Ошибка в запросе"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id} [put]
func UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Println("INFO: Failed to decode request body for update")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
		return
	}

	var existingSong models.Song
	if err := database.DB.First(&existingSong, id).Error; err != nil {
		log.Println("INFO: Song not found with ID:", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Song not found",
		})
		return
	}

	if err := database.DB.Model(&existingSong).Updates(song).Error; err != nil {
		log.Println("INFO: Failed to update song with ID:", id)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update song",
		})
		return
	}

	log.Println("DEBUG: Successfully updated song with ID:", id)
	w.WriteHeader(http.StatusNoContent)
}

// AddSong добавляет новую песню в базу данных.
// @Summary Добавление новой песни
// @Description Добавляет новую песню в базу данных, сначала запрашивая информацию из внешнего API.
// @Accept json
// @Produce json
// @Param song body models.Song true "Данные о песне"
// @Success 201 {object} models.Song "Добавленная песня"
// @Failure 400 {object} models.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs [post]
func AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Println("INFO: Failed to decode request body for adding song")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
		return
	}

	if song.Group == "" || song.Name == "" {
		log.Println("INFO: Group or song name is empty")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Group and song name must not be empty",
		})
		return
	}

	// мой тестовый Mock API запускается на порту 8081
	apiURL := "http://localhost:8081/info?group=" + url.QueryEscape(song.Group) + "&song=" + url.QueryEscape(song.Name)
	resp, err := http.Get(apiURL)

	if err != nil {
		log.Println("INFO: Error calling external API:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to call external API",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("INFO: Non-200 response from external API:", resp.StatusCode)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve data from external API",
		})
		return
	}

	var externalSongDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&externalSongDetail); err != nil {
		log.Println("INFO: Failed to decode external API response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to decode external API response",
		})
		return
	}

	song.ReleaseDate = externalSongDetail.ReleaseDate
	song.Text = externalSongDetail.Text

	song.Link = externalSongDetail.Link

	if err := database.DB.Create(&song).Error; err != nil {
		log.Println("INFO: Failed to save song to the database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to save song to the database",
		})
		return
	}

	log.Println("DEBUG: Successfully added song")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}
