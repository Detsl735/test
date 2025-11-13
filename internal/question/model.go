package question

import (
	"time"
)

type Question struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Text      string    `gorm:"type:text" json:"text"`
	CreatedAt time.Time `gorm:"type:autoCreateTime" json:"created_at"`
}

type CreateQuestionRequest struct {
	Text string `json:"text" validate:"required"`
}
