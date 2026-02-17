package models

import "time"

type PlaybackProgress struct {
	UserID    string    `json:"userId"`
	ContentID string    `json:"contentId"`
	Seconds   int       `json:"seconds"` // tiempo reproducido
	Percent   float64   `json:"percent"` // 0..100
	Completed bool      `json:"completed"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}
