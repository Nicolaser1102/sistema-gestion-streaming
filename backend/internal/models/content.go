package models

import "time"

type ContentType string

const (
	ContentMovie  ContentType = "MOVIE"
	ContentSeries ContentType = "SERIES"
)

type Content struct {
	ID          string      `json:"id"`
	Type        ContentType `json:"type"`
	Title       string      `json:"title"`
	Synopsis    string      `json:"synopsis"`
	Genre       string      `json:"genre"`
	Year        int         `json:"year"`
	DurationMin int         `json:"durationMin"`
	CoverURL    string      `json:"coverURL"`
	Active      bool        `json:"active"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}
