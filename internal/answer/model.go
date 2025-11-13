package answer

import (
	"time"
)

type Answer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	QuestionID uint      `gorm:"not null;index" json:"question_id"`
	UserID     string    `gorm:"type:varchar(64);not null;" json:"user_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type CreateAnswerRequest struct {
	QuestionID uint   `json:"question_id" validate:"required"`
	UserID     string `json:"user_id" validate:"required"`
	Text       string `json:"text" validate:"required"`
}
