package models

import "time"

type MyListItem struct {
	UserID    string    `json:"userId"`
	ContentID string    `json:"contentId"`
	CreatedAt time.Time `json:"createdAt"`
}
