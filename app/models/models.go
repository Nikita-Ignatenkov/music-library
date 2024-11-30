package models

import (
	"time"
)

// MyBaseModel добавляет стандартные поля для других моделей.
type MyBaseModel struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// Song представляет песню в системе.
// @Description Структура песни
type Song struct {
	MyBaseModel        // Включает поля ID, CreatedAt, UpdatedAt и DeletedAt
	Group       string `json:"group" gorm:"column:artist;not null"` // Группа или исполнитель
	Name        string `json:"song" gorm:"not null"`                // Название песни
	ReleaseDate string `json:"releaseDate"`                         // Дата релиза
	Text        string `json:"text"`                                // Текст песни
	Link        string `json:"link"`                                // Ссылка на песню
}

// SongDetail содержит дополнительные детали о песне.
// @Description Структура с деталями песни
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"` // Дата релиза
	Text        string `json:"text"`        // Текст песни
	Link        string `json:"link"`        // Ссылка на песню
}

// ErrorResponse описывает ответ с ошибкой.
// @Description Структура ответа для ошибок API.
type ErrorResponse struct {
	Code    int    `json:"code"`    // Код ошибки
	Message string `json:"message"` // Сообщение об ошибке
}
