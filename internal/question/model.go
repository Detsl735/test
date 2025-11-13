package question

import (
	"testTask/internal/answer"
	"time"
)

type Question struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Text      string    `gorm:"type:text; not null" json:"text"`
	CreatedAt time.Time `gorm:"type:autoCreateTime" json:"created_at"`

	Answers []answer.Answer `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}

type CreateQuestionRequest struct {
	Text string `json:"text" validate:"required"`
}
