package domain

import "time"

type SessionState string

const (
	StateNotStarted  SessionState = "not_started"
	StateActive      SessionState = "active"
	StateSubmitted   SessionState = "submitted"
	StateExpired     SessionState = "expired"
	StateUnderReview SessionState = "under_review"
)

type QuizSession struct {
	ID              uint         `gorm:"primaryKey" json:"id"`
	QuizID          uint         `gorm:"uniqueIndex:idx_session_quiz_student;not null" json:"quiz_id"`
	StudentID       uint         `gorm:"uniqueIndex:idx_session_quiz_student;not null" json:"student_id"`
	StartTime       *time.Time   `json:"start_time"` // Nullable if not started
	State           SessionState `gorm:"not null;default:'not_started'" json:"state"`
	SuspiciousScore int          `gorm:"not null;default:0" json:"suspicious_score"`
	TabSwitchCount  int          `gorm:"not null;default:0" json:"tab_switch_count"`
	Score           float64      `gorm:"not null;default:0" json:"score"` // Final calculated score
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`

	SessionQuestions []SessionQuestion `gorm:"foreignKey:SessionID" json:"session_questions,omitempty"`
}

type SessionQuestion struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	SessionID  uint `gorm:"index;not null" json:"session_id"`
	QuestionID uint `gorm:"not null" json:"question_id"`
	OrderIndex int  `gorm:"not null" json:"order_index"` // For shuffling

	SessionOptions []SessionOption `gorm:"foreignKey:SessionQuestionID" json:"session_options,omitempty"`
	Answer         *SessionAnswer  `gorm:"foreignKey:SessionQuestionID" json:"answer,omitempty"`
}

type SessionOption struct {
	ID                uint `gorm:"primaryKey" json:"id"`
	SessionQuestionID uint `gorm:"index;not null" json:"session_question_id"`
	OptionID          uint `gorm:"not null" json:"option_id"`
	OrderIndex        int  `gorm:"not null" json:"order_index"` // For shuffling
}

type SessionAnswer struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SessionQuestionID uint      `gorm:"uniqueIndex;not null" json:"session_question_id"`
	SelectedOptionID  uint      `gorm:"not null" json:"selected_option_id"`
	CreatedAt         time.Time `json:"created_at"`
}
