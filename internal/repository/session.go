package repository

import (
	"context"
	"online-quiz/internal/domain"

	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.QuizSession) error
	GetSessionByQuizAndStudent(ctx context.Context, quizID, studentID uint) (*domain.QuizSession, error)
	GetSessionByID(ctx context.Context, id uint) (*domain.QuizSession, error)
	UpdateSession(ctx context.Context, session *domain.QuizSession) error
	CreateSessionAnswers(ctx context.Context, answers []domain.SessionAnswer) error
	WithTransaction(ctx context.Context, fn func(txRepo SessionRepository) error) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db}
}

func (r *sessionRepository) WithTransaction(ctx context.Context, fn func(txRepo SessionRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := NewSessionRepository(tx)
		return fn(txRepo)
	})
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *domain.QuizSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) GetSessionByQuizAndStudent(ctx context.Context, quizID, studentID uint) (*domain.QuizSession, error) {
	var session domain.QuizSession
	if err := r.db.WithContext(ctx).Preload("SessionQuestions.SessionOptions").
		Where("quiz_id = ? AND student_id = ?", quizID, studentID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id uint) (*domain.QuizSession, error) {
	var session domain.QuizSession
	if err := r.db.WithContext(ctx).Preload("SessionQuestions.SessionOptions").First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) UpdateSession(ctx context.Context, session *domain.QuizSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *sessionRepository) CreateSessionAnswers(ctx context.Context, answers []domain.SessionAnswer) error {
	return r.db.WithContext(ctx).Create(&answers).Error
}
