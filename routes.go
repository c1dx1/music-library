package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"net/http"
	"strconv"
	"strings"
)

func searchFreeArgNum(m map[int]int) int {
	for i := 1; i <= len(m); i++ {
		if m[i] == 0 {
			m[i] = 1
			return i
		}
	}
	return -1
}

// @Summary Get songs with optional filters and pagination
// @Description Получение данных библиотеки с фильтрацией по всем полям и пагинацией
// @Tags Songs
// @Param id query int false "ID песни"
// @Param group query string false "Название группы"
// @Param name query string false "Название песни"
// @Param releaseDate query string false "Дата релиза"
// @Param text query string false "Текст песни"
// @Param link query string false "Ссылка"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Produce  json
// @Success 200 {array} Song
// @Failure 500 {string} string "Internal server error"
// @Router /songs [get]
func getSongs(w http.ResponseWriter, r *http.Request) {
	log.Debug("Entering getSongs function")

	idStr := r.URL.Query().Get("id")
	group := r.URL.Query().Get("group")
	name := r.URL.Query().Get("name")
	releaseDate := r.URL.Query().Get("releaseDate")
	text := r.URL.Query().Get("text")
	link := r.URL.Query().Get("link")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	log.Debugf("Query parameters: id=%s, group=%s, name=%s, releaseDate=%s, text=%s, link=%s, page=%s, limit=%s",
		idStr, group, name, releaseDate, text, link, pageStr, limitStr)

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
		log.Info("Page number is less than 1, defaulting to 1")
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
		log.Info("Limit is less than 1, defaulting to 10")
	}

	offset := (page - 1) * limit
	log.Debugf("Calculating offset: %d", offset)

	query := "SELECT id, group_name, song_name, release_date, text, link FROM musicdb WHERE 1=1"
	args := []interface{}{}

	m := map[int]int{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
		6: 0,
		7: 0,
		8: 0,
	}

	if idStr != "" {
		id, _ := strconv.Atoi(idStr)
		query += fmt.Sprintf(" AND id=$%d", searchFreeArgNum(m))
		args = append(args, id)
		log.Debugf("Filtering by id: %d", id)
	}
	if group != "" {
		query += fmt.Sprintf(" AND group_name ILIKE $%d", searchFreeArgNum(m))
		args = append(args, "%"+group+"%")
		log.Debugf("Filtering by group: %s", group)
	}
	if name != "" {
		query += fmt.Sprintf(" AND song_name ILIKE $%d", searchFreeArgNum(m))
		args = append(args, "%"+name+"%")
		log.Debugf("Filtering by song name: %s", name)
	}
	if releaseDate != "" {
		query += fmt.Sprintf(" AND release_date ILIKE $%d", searchFreeArgNum(m))
		args = append(args, "%"+releaseDate+"%")
		log.Debugf("Filtering by release date: %s", releaseDate)
	}
	if text != "" {
		query += fmt.Sprintf(" AND text ILIKE $%d", searchFreeArgNum(m))
		args = append(args, "%"+text+"%")
		log.Debugf("Filtering by text: %s", text)
	}
	if link != "" {
		query += fmt.Sprintf(" AND link ILIKE $%d", searchFreeArgNum(m))
		args = append(args, "%"+link+"%")
		log.Debugf("Filtering by link: %s", link)
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", searchFreeArgNum(m), searchFreeArgNum(m))
	args = append(args, limit, offset)

	log.Debugf("Executing query: %s with args: %+v", query, args)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Errorf("Database query error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		err := rows.Scan(&song.ID, &song.Group, &song.Name, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			log.Errorf("Row scan error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		songs = append(songs, song)
		log.Debugf("Retrieved song: %+v", song)
	}

	log.Infof("Returning %d songs", len(songs))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// @Summary Get song lyrics with pagination
// @Description Получение текста песни по ID с пагинацией по куплетам
// @Tags Songs
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество куплетов на странице (по умолчанию 3)"
// @Produce  json
// @Success 200 {array} string
// @Failure 500 {string} string "Internal server error"
// @Router /songs/{id} [get]
func getTextFromSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Debugf("Fetching text for song with id: %s", id)
	query := "SELECT id, group_name, song_name, release_date, text, link FROM musicdb WHERE id=$1"

	var song Song
	err := db.QueryRow(query, id).Scan(&song.ID, &song.Group, &song.Name, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		log.Errorf("Error fetching song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	verses := paginateVerses(song.Text, r)
	log.Infof("Returning %d verses for song with id: %s", len(verses), id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(verses)
}

func paginateVerses(text string, r *http.Request) []string {
	verses := strings.Split(text, "\n\n")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 3
	}

	start := (page - 1) * limit
	end := start + limit
	if start > len(verses) {
		return []string{}
	}
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end]
}

// @Summary Add a new song
// @Description Добавление новой песни с запросом данных о песне из внешнего API
// @Tags Songs
// @Param group query string true "Название группы"
// @Param name query string true "Название песни"
// @Produce  json
// @Success 201 {object} Song
// @Failure 500 {string} string "Internal server error"
// @Router /songs [post]
func addSong(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	name := r.URL.Query().Get("name")

	q := fmt.Sprintf("%s %s", group, name)
	log.Debugf("Adding song with group: %s, name: %s", group, name)

	song, err := getGeniusData(q)
	if err != nil {
		log.Errorf("Error retrieving song data from Genius: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO musicdb (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5) returning id"
	err = db.QueryRow(query, song.Group, song.Name, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID)
	if err != nil {
		log.Errorf("Error inserting song into database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Song added with id: %d", song.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// @Summary Delete a song
// @Description Удаление песни по ID
// @Tags Songs
// @Param id path int true "ID песни"
// @Success 204 "No Content"
// @Failure 500 {string} string "Internal server error"
// @Router /songs/{id} [delete]
func deleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Debugf("Deleting song with id: %s", id)
	query := "DELETE FROM musicdb WHERE id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		log.Errorf("Error deleting song from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Song with id %s deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Edit an existing song
// @Description Редактирование данных песни по ID, обновляются только переданные поля
// @Tags Songs
// @Param id path int true "ID песни"
// @Param song body Song true "Данные для обновления песни"
// @Produce  json
// @Success 204 "No Content"
// @Failure 500 {string} string "Internal server error"
// @Router /songs/{id} [put]
func editSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var currentSong Song
	err := db.QueryRow("SELECT group_name, song_name, release_date, text, link FROM musicdb WHERE id = $1", id).Scan(
		&currentSong.Group, &currentSong.Name, &currentSong.ReleaseDate, &currentSong.Text, &currentSong.Link,
	)
	if err != nil {
		log.Errorf("Error fetching current song from database: %v", err)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	var inputSong Song
	err = json.NewDecoder(r.Body).Decode(&inputSong)
	log.Debugf("Song: %s", inputSong)
	if err != nil {
		log.Warnf("Invalid input for inputSong editing: %s", err.Error())
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if inputSong.Group == "" {
		inputSong.Group = currentSong.Group
	}
	if inputSong.Name == "" {
		inputSong.Name = currentSong.Name
	}
	if inputSong.ReleaseDate.Year() == 0001 {
		inputSong.ReleaseDate = currentSong.ReleaseDate
	}
	if inputSong.Text == "" {
		inputSong.Text = currentSong.Text
	}
	if inputSong.Link == "" {
		inputSong.Link = currentSong.Link
	}

	log.Debugf("Editing song with id: %s", id)
	query := "UPDATE musicdb SET group_name=$1, song_name=$2, release_date=$3, text=$4, link=$5 WHERE id=$6"
	_, err = db.Exec(query, inputSong.Group, inputSong.Name, inputSong.ReleaseDate, inputSong.Text, inputSong.Link, id)
	if err != nil {
		log.Errorf("Error updating song in database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Song with id %s updated", id)
	w.WriteHeader(http.StatusNoContent)
}
