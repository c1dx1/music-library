package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type GeniusSearch struct {
	Hits []struct {
		Result struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Artists string `json:"primary_artist_names"`
		} `json:"result"`
	} `json:"hits"`
}

type GeniusLyrics struct {
	Lyrics struct {
		Lyrics struct {
			Body struct {
				HTML string `json:"html"`
			} `json:"body"`
		} `json:"lyrics"`
	} `json:"lyrics"`
}

type GeniusDetail struct {
	Song struct {
		ReleaseDate string `json:"release_date"`
		Link        string `json:"youtube_url"`
	}
}

func getGeniusSearchData(q string) (*GeniusSearch, error) {
	log.Debugf("Fetching Genius search data for query: %s", q)
	url := fmt.Sprintf(os.Getenv("URL_GENIUS_SEARCH"), q)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-key", os.Getenv("RAPIDAPI_KEY"))
	req.Header.Add("x-rapidapi-host", os.Getenv("RAPIDAPI_HOST"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error making request to Genius API: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	var geniusIDReleaseDate GeniusSearch
	if err := json.NewDecoder(res.Body).Decode(&geniusIDReleaseDate); err != nil {
		log.Errorf("Error decoding response from Genius API: %v", err)
		return nil, err
	}

	log.Infof("Successfully fetched Genius search data")
	return &geniusIDReleaseDate, nil
}

func getLyrics(id int) (*GeniusLyrics, error) {
	log.Debugf("Fetching lyrics for song ID: %d", id)
	url := fmt.Sprintf(os.Getenv("URL_GENIUS_LYRICS"), id)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-key", os.Getenv("RAPIDAPI_KEY"))
	req.Header.Add("x-rapidapi-host", os.Getenv("RAPIDAPI_HOST"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error making request to Genius API for lyrics: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	var geniusLyrics GeniusLyrics
	if err := json.NewDecoder(res.Body).Decode(&geniusLyrics); err != nil {
		log.Errorf("Error decoding response from Genius API for lyrics: %v", err)
		return nil, err
	}

	log.Infof("Successfully fetched lyrics for song ID: %d", id)
	return &geniusLyrics, nil
}

func getDetails(id int) (*GeniusDetail, error) {
	log.Debugf("Fetching details for song ID: %d", id)
	url := fmt.Sprintf(os.Getenv("URL_GENIUS_DETAILS"), id)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-key", os.Getenv("RAPIDAPI_KEY"))
	req.Header.Add("x-rapidapi-host", os.Getenv("RAPIDAPI_HOST"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error making request to Genius API for details: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	var geniusDetail GeniusDetail
	if err := json.NewDecoder(res.Body).Decode(&geniusDetail); err != nil {
		log.Errorf("Error decoding response from Genius API for details: %v", err)
		return nil, err
	}

	log.Infof("Successfully fetched details for song ID: %d", id)
	return &geniusDetail, nil
}

func formatLyrics(lyrics string) string {
	lyrics = strings.ReplaceAll(lyrics, "<br><br>", "\\n")

	re := regexp.MustCompile(`</?[^>]+>`)
	formattedLyrics := re.ReplaceAllString(lyrics, "")

	return formattedLyrics
}

func getGeniusData(q string) (*Song, error) {
	log.Debugf("Fetching Genius data for query: %s", q)
	var song Song

	geniusSearch, err := getGeniusSearchData(q)
	if err != nil {
		return nil, err
	}

	song.Group = geniusSearch.Hits[0].Result.Artists
	song.Name = geniusSearch.Hits[0].Result.Title

	id := geniusSearch.Hits[0].Result.ID

	geniusLyricsHTML, err := getLyrics(id)
	if err != nil {
		return nil, err
	}

	lyrics := formatLyrics(geniusLyricsHTML.Lyrics.Lyrics.Body.HTML)
	song.Text = lyrics

	geniusDetail, err := getDetails(id)
	if err != nil {
		return nil, err
	}

	releaseDate, _ := time.Parse("2006-01-02", geniusDetail.Song.ReleaseDate)
	song.ReleaseDate = releaseDate
	song.Link = geniusDetail.Song.Link

	log.Infof("Successfully fetched Genius data for song: %s by %s", song.Name, song.Group)
	return &song, nil
}
