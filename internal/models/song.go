package models

import (
	"encoding/json"
	"fmt"
	"time"
)

//type Song struct {
//	ID          *int       `json:"id" example:"1"`
//	Group       *string    `json:"group" example:"The Beatles"`
//	Song        *string    `json:"song" example:"Hey Jude"`
//	ReleaseDate *time.Time `json:"releaseDate" time_format:"2006-02-01" example:"1968-08-26"`
//	Text        *string    `json:"text" example:"Hey, Jude, don't make it bad\nTake a sad song and make it better\nRemember to let her into your heart\nThen you can start to make it better\n\nHey, Jude, don't be afraid You were made to go out and get her\nThe minute you let her under your skin\nThen you begin to make it better\nAnd anytime you feel the pain, hey, Jude, refrain\nDon't carry the world upon your shoulders\nFor well you know that it's a fool who plays it cool\nBy making his world a little colder\nNa-na-na-na-na, na-na-na-na\n\nHey, Jude, don't let me down\n\nYou have found her, now go and get her\n(Let it out and let it in)\nRemember (Hey, Jude) to let her into your heart\nThen you can start to make it better"`
//	Link        *string    `json:"link" example:"https://example.com/heyjude"`
//}

type Song struct {
	ID          *int       `json:"id" example:"1"`
	Group       *string    `json:"group" example:"The Beatles"`
	Song        *string    `json:"song" example:"Hey Jude"`
	ReleaseDate *time.Time `json:"releaseDate" time_format:"2006-01-02" example:"1968-08-26"`
	Text        *string    `json:"text" example:"Hey, Jude, don't make it bad\nTake a sad song and make it better\nRemember to let her into your heart\nThen you can start to make it better\n\nHey, Jude, don't be afraid You were made to go out and get her\nThe minute you let her under your skin\nThen you begin to make it better\nAnd anytime you feel the pain, hey, Jude, refrain\nDon't carry the world upon your shoulders\nFor well you know that it's a fool who plays it cool\nBy making his world a little colder\nNa-na-na-na-na, na-na-na-na\n\nHey, Jude, don't let me down\n\nYou have found her, now go and get her\n(Let it out and let it in)\nRemember (Hey, Jude) to let her into your heart\nThen you can start to make it better"`
	Link        *string    `json:"link" example:"https://example.com/heyjude"`
}

func (s *Song) MarshalJSON() ([]byte, error) {
	type Alias struct {
		ID          *int    `json:"id"`
		Group       *string `json:"group"`
		Song        *string `json:"song"`
		ReleaseDate *string `json:"releaseDate"`
		Text        *string `json:"text"`
		Link        *string `json:"link"`
	}

	var formattedDate *string
	if s.ReleaseDate != nil {
		formattedDateStr := s.ReleaseDate.Format("2006-01-02")
		formattedDate = &formattedDateStr
	}

	aux := &Alias{
		ID:          s.ID,
		Group:       s.Group,
		Song:        s.Song,
		ReleaseDate: formattedDate,
		Text:        s.Text,
		Link:        s.Link,
	}
	return json.Marshal(aux)
}

func (s *Song) UnmarshalJSON(data []byte) error {
	type Alias Song

	aux := &struct {
		RawDate *string `json:"releaseDate"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.RawDate != nil {
		const layoutDefault = "2006-01-02"
		const layoutExternal = "02.01.2006"
		var parsedDate time.Time
		if *aux.RawDate == "" {
			var err error
			parsedDate, err = time.Parse(layoutDefault, "0001-01-01")
			if err != nil {
				return fmt.Errorf("invalid date format: %w", err)
			}
		} else {
			var err error
			parsedDate, err = time.Parse(layoutDefault, *aux.RawDate)
			if err != nil {
				parsedDate, err = time.Parse(layoutExternal, *aux.RawDate)
				if err != nil {
					return fmt.Errorf("invalid date format: %w", err)
				}
			}
		}
		s.ReleaseDate = &parsedDate
	}
	return nil
}

//type EditSong struct {
//	Group   string        `json:"group" example:"The Beatles"`
//	Song    string    	  `json:"song" example:"Hey Jude"`
//	ReleaseDate time.Time `json:"releaseDate" time_format:"2006-01-02" example:"1968-08-26"`
//	Text        string    `json:"text" example:"Hey, Jude, don't make it bad\nTake a sad song and make it better\nRemember to let her into your heart\nThen you can start to make it better\n\nHey, Jude, don't be afraid You were made to go out and get her\nThe minute you let her under your skin\nThen you begin to make it better\nAnd anytime you feel the pain, hey, Jude, refrain\nDon't carry the world upon your shoulders\nFor well you know that it's a fool who plays it cool\nBy making his world a little colder\nNa-na-na-na-na, na-na-na-na\n\nHey, Jude, don't let me down\n\nYou have found her, now go and get her\n(Let it out and let it in)\nRemember (Hey, Jude) to let her into your heart\nThen you can start to make it better"`
//	Link        string    `json:"link" example:"https://example.com/heyjude"`
//}

type EditSong struct {
	Group       *string    `json:"group" example:"The Beatles"`
	Song        *string    `json:"song" example:"Hey Jude"`
	ReleaseDate *time.Time `json:"releaseDate" time_format:"2006-01-02" example:"1968-08-26"`
	Text        *string    `json:"text" example:"Hey, Jude, don't make it bad\nTake a sad song and make it better\nRemember to let her into your heart\nThen you can start to make it better\n\nHey, Jude, don't be afraid You were made to go out and get her\nThe minute you let her under your skin\nThen you begin to make it better\nAnd anytime you feel the pain, hey, Jude, refrain\nDon't carry the world upon your shoulders\nFor well you know that it's a fool who plays it cool\nBy making his world a little colder\nNa-na-na-na-na, na-na-na-na\n\nHey, Jude, don't let me down\n\nYou have found her, now go and get her\n(Let it out and let it in)\nRemember (Hey, Jude) to let her into your heart\nThen you can start to make it better"`
	Link        *string    `json:"link" example:"https://example.com/heyjude"`
}

func (s *EditSong) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Group       *string `json:"group"`
		Song        *string `json:"song"`
		ReleaseDate *string `json:"releaseDate"`
		Text        *string `json:"text"`
		Link        *string `json:"link"`
	}

	var formattedDate *string
	if s.ReleaseDate != nil {
		formattedDateStr := s.ReleaseDate.Format("2006-01-02")
		formattedDate = &formattedDateStr
	}

	aux := &Alias{
		Group:       s.Group,
		Song:        s.Song,
		ReleaseDate: formattedDate,
		Text:        s.Text,
		Link:        s.Link,
	}
	return json.Marshal(aux)
}

func (s *EditSong) UnmarshalJSON(data []byte) error {
	type Alias EditSong

	aux := &struct {
		RawDate *string `json:"releaseDate"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.RawDate != nil {
		const layoutDefault = "2006-01-02"
		const layoutExternal = "02.01.2006"
		var parsedDate time.Time
		if *aux.RawDate == "" {
			var err error
			parsedDate, err = time.Parse(layoutDefault, "0001-01-01")
			if err != nil {
				return fmt.Errorf("invalid date format: %w", err)
			}
		} else {
			var err error
			parsedDate, err = time.Parse(layoutDefault, *aux.RawDate)
			if err != nil {
				parsedDate, err = time.Parse(layoutExternal, *aux.RawDate)
				if err != nil {
					return fmt.Errorf("invalid date format: %w", err)
				}
			}
		}
		s.ReleaseDate = &parsedDate
	}
	return nil
}

//type AddSong struct {
//	Group       *string    `json:"group" example:"The Beatles"`
//	Song        *string    `json:"song" example:"Hey Jude"`
//}

type AddSong struct {
	Group *string `json:"group" example:"The Beatles"`
	Song  *string `json:"song" example:"Hey Jude"`
}
