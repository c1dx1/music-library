package models

import "time"

type LibraryFilter struct {
	ID          *int       `form:"id" example:"1"`
	Group       *string    `form:"group" example:"The Beatles"`
	Song        *string    `form:"song" example:"Hey Jude"`
	ReleaseDate *time.Time `form:"releaseDate" time_format:"2006-01-02" example:"1968-08-26"`
	Text        *string    `form:"text" example:"Hey, Jude, don't make it bad\nTake a sad song and make it better\nRemember to let her into your heart\nThen you can start to make it better\n\nHey, Jude, don't be afraid You were made to go out and get her\nThe minute you let her under your skin\nThen you begin to make it better\nAnd anytime you feel the pain, hey, Jude, refrain\nDon't carry the world upon your shoulders\nFor well you know that it's a fool who plays it cool\nBy making his world a little colder\nNa-na-na-na-na, na-na-na-na\n\nHey, Jude, don't let me down\n\nYou have found her, now go and get her\n(Let it out and let it in)\nRemember (Hey, Jude) to let her into your heart\nThen you can start to make it better"`
	Link        *string    `form:"link" example:"https://example.com/heyjude"`
	Page        *int       `form:"page" example:"1"`
	Limit       *int       `form:"limit" example:"10"`
}
