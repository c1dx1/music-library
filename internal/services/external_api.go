package services

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"music-library/config"
	"music-library/internal/models"
	"net/http"
	"strings"
)

type ExternalAPIClient struct {
	cfg *config.Config
	log *logrus.Logger
}

func NewExternalAPIClient(cfg *config.Config, log *logrus.Logger) *ExternalAPIClient {
	return &ExternalAPIClient{cfg: cfg, log: log}
}

func (e *ExternalAPIClient) GetSongDetails(song *models.Song) error {
	e.log.Infof("GetSongDetails called for Group: %s, Song: %s", *song.Group, *song.Song)

	groupWithoutSpaces := strings.ReplaceAll(*song.Group, " ", "%20")
	songWithoutSpaces := strings.ReplaceAll(*song.Song, " ", "%20")

	url := fmt.Sprintf("%s/info?group=%s&song=%s", e.cfg.ExternalAPI, groupWithoutSpaces, songWithoutSpaces)
	e.log.Debugf("Making GET request to URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		e.log.Errorf("HTTP GET request failed: %v", err)
		return fmt.Errorf("external_api: getSongDetails: http get: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			e.log.Warnf("Failed to close response body: %v", cerr)
		}
	}()
	e.log.Infof("Received response with status code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			e.log.Warn("External API returned BadRequest (400)")
			return fmt.Errorf("external_api: getSongDetails: BadRequest")
		} else {
			e.log.Error("External API returned an unexpected status code")
			return fmt.Errorf("external_api: getSongDetails: InternalServerError")
		}
	}

	if err = json.NewDecoder(resp.Body).Decode(&song); err != nil {
		e.log.Errorf("Failed to decode JSON response: %v", err)
		return fmt.Errorf("external_api: getSongDetails: json decode: %w", err)
	}
	e.log.Infof("Successfully decoded song details: %+v", song)

	return nil
}
