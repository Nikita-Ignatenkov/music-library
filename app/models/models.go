package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Group       string `json:"group" gorm:"not null"`
	Name        string `json:"song" gorm:"not null"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
