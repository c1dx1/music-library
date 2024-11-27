package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"music-library/internal/models"
	"music-library/internal/services"
	"net/http"
	"strconv"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MLibHandler struct {
	Service *services.MLibService
	log     *logrus.Logger
}

func NewMLibHandler(service *services.MLibService, log *logrus.Logger) *MLibHandler {
	return &MLibHandler{Service: service, log: log}
}

// GetLibrary godoc
// @Summary      Get the music library
// @Description  Fetches the music library with optional filters
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id 		 query    int    false "Filter by id"
// @Param        group       query    string false "Filter by group"
// @Param        song      	 query    string false "Filter by song"
// @Param        releaseDate query    string false "Filter by release date"
// @Param        text      	 query    string false "Filter by text"
// @Param        link      	 query    string false "Filter by link"
// @Param        page      	 query    int    false "Page number" default(1)
// @Param        limit       query    int    false "Page size" default(10)
// @Success      200         {array}  models.Song
// @Failure      400         {object} ErrorResponse
// @Failure      500         {object} ErrorResponse
// @Router       /songs [get]
func (h *MLibHandler) GetLibrary(c *gin.Context) {
	h.log.Debug("Entering GetLibrary handler")

	var filter models.LibraryFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.log.Warnf("Failed to bind query parameters: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Debugf("Fetching library with filter: %+v", filter)
	songs, err := h.Service.GetLibrary(filter)
	if err != nil {
		h.log.Errorf("Failed to fetch library: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Info("Successfully fetched library")
	c.JSON(http.StatusOK, songs)
}

// GetText godoc
// @Summary      Get song text
// @Description  Fetches the text of a song by ID with pagination
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id          path     int     true  "Song ID"
// @Param        page        query    int     false "Page number" default(1)
// @Param        limit       query    int     false "Page size"   default(10)
// @Success      200         {array}  string
// @Failure      400         {object} ErrorResponse
// @Failure      500         {object} ErrorResponse
// @Router       /songs/{id} [get]
func (h *MLibHandler) GetText(c *gin.Context) {
	h.log.Debug("Entering GetText handler")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warnf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.log.Warnf("Invalid page parameter: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.log.Warnf("Invalid limit parameter: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Debugf("Fetching text for song ID %d with page %d and limit %d", id, page, limit)
	verses, err := h.Service.GetText(id, page, limit)
	if err != nil {
		h.log.Errorf("Failed to fetch text: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Info("Successfully fetched text for song")
	c.JSON(http.StatusOK, verses)
}

// DeleteSong godoc
// @Summary      Delete a song
// @Description  Deletes a song by ID
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id          path     int     true  "Song ID"
// @Success      200         {object} SuccessResponse
// @Failure      400         {object} ErrorResponse
// @Failure      500         {object} ErrorResponse
// @Router       /songs/{id} [delete]
func (h *MLibHandler) DeleteSong(c *gin.Context) {
	h.log.Debug("Entering DeleteSong handler")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warnf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Debugf("Attempting to delete song with ID %d", id)
	err = h.Service.DeleteSong(id)
	if err != nil {
		h.log.Errorf("Failed to delete song: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Info("Successfully deleted song")
	c.JSON(http.StatusOK, SuccessResponse{Message: "Song deleted successfully"})
}

// EditSong godoc
// @Summary      Edit a song
// @Description  Edits a song's details by ID
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id          path     int       true  "Song ID"
// @Param        song        body     models.EditSong true "Updated song data"
// @Success      200         {object} SuccessResponse
// @Failure      400         {object} ErrorResponse
// @Failure      500         {object} ErrorResponse
// @Router       /songs/{id} [put]
func (h *MLibHandler) EditSong(c *gin.Context) {
	h.log.Debug("Entering EditSong handler")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warnf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var editSong models.EditSong
	if err := c.ShouldBindJSON(&editSong); err != nil {
		h.log.Warnf("Failed to bind JSON for edit request: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Debugf("Editing song with ID %d: %+v", id, editSong)
	err = h.Service.EditSong(id, editSong)
	if err != nil {
		h.log.Errorf("Failed to edit song: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Info("Successfully edited song")
	c.JSON(http.StatusOK, SuccessResponse{Message: "Song edited successfully"})
}

// AddSong godoc
// @Summary      Add a new song
// @Description  Adds a new song to the library
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song        body     models.AddSong true "New song data"
// @Success      200         {object} SuccessResponse
// @Failure      400         {object} ErrorResponse
// @Failure      500         {object} ErrorResponse
// @Router       /songs [post]
func (h *MLibHandler) AddSong(c *gin.Context) {
	h.log.Debug("Entering AddSong handler")

	var addSong models.AddSong
	if err := c.ShouldBindJSON(&addSong); err != nil {
		h.log.Warnf("Failed to bind JSON for new song: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	song := models.Song{Group: addSong.Group, Song: addSong.Song}

	h.log.Debugf("Adding new song: %+v", song)
	err := h.Service.AddSong(song)
	if err != nil {
		h.log.Errorf("Failed to add song: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Info("Successfully added new song")
	c.JSON(http.StatusOK, SuccessResponse{Message: "Song added successfully"})
}
