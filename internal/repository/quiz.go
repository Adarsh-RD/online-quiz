package repository

import (
	"context"
	"online-quiz/internal/domain"

	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(ctx context.Context, quiz *domain.Quiz) error
	GetByID(ctx context.Context, id uint) (*domain.Quiz, error)
	GetAll(ctx context.Context) ([]domain.Quiz, error)
	GetPublished(ctx context.Context) ([]domain.Quiz, error)
	CreateQuestion(ctx context.Context, question *domain.Question) error
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db}
}

func (r *quizRepository) Create(ctx context.Context, quiz *domain.Quiz) error {
	return r.db.WithContext(ctx).Create(quiz).Error
}

func (r *quizRepository) GetByID(ctx context.Context, id uint) (*domain.Quiz, error) {
	var quiz domain.Quiz
	if err := r.db.WithContext(ctx).Preload("Questions.Options").First(&quiz, id).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) GetAll(ctx context.Context) ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.WithContext(ctx).Find(&quizzes).Error
	return quizzes, err
}

func (r *quizRepository) GetPublished(ctx context.Context) ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.WithContext(ctx).Where("published = ?", true).Find(&quizzes).Error
	return quizzes, err
}

func (r *quizRepository) CreateQuestion(ctx context.Context, question *domain.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}
