package domain

import "time"

type Quiz struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Title     string     `gorm:"not null" json:"title"`
	StartTime time.Time  `gorm:"not null" json:"start_time"`
	EndTime   time.Time  `gorm:"not null" json:"end_time"`
	TeacherID uint       `gorm:"index;not null" json:"teacher_id"`
	Published bool       `gorm:"not null;default:false" json:"published"`
	Questions []Question `gorm:"foreignKey:QuizID" json:"questions,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Question struct {
	ID       uint     `gorm:"primaryKey" json:"id"`
	QuizID   uint     `gorm:"index;not null" json:"quiz_id"`
	Text     string   `gorm:"not null" json:"text"`
	Marks    float64  `gorm:"not null;default:1" json:"marks"`
	Options  []Option `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
}

type Option struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	QuestionID uint   `gorm:"index;not null" json:"question_id"`
	Text       string `gorm:"not null" json:"text"`
	IsCorrect  bool   `gorm:"not null" json:"is_correct"`
}
