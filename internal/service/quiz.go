package service

import (
	"context"
	"errors"
	"time"

	"online-quiz/internal/domain"
	"online-quiz/internal/repository"
)

type QuizService interface {
	CreateQuiz(ctx context.Context, quiz *domain.Quiz) error
	GetQuizzesForStudent(ctx context.Context) ([]domain.Quiz, error)
	GetQuizForTeacher(ctx context.Context, id uint) (*domain.Quiz, error)
}

type quizService struct {
	quizRepo repository.QuizRepository
}

func NewQuizService(quizRepo repository.QuizRepository) QuizService {
	return &quizService{quizRepo}
}

func (s *quizService) CreateQuiz(ctx context.Context, quiz *domain.Quiz) error {
	if quiz.EndTime.Before(quiz.StartTime) {
		return errors.New("end time cannot be before start time")
	}
	return s.quizRepo.Create(ctx, quiz)
}

func (s *quizService) GetQuizzesForStudent(ctx context.Context) ([]domain.Quiz, error) {
	// Only return published quizzes
	quizzes, err := s.quizRepo.GetPublished(ctx)
	if err != nil {
		return nil, err
	}
	
	now := time.Now()
	var availableQuizzes []domain.Quiz
	for _, q := range quizzes {
		// Only show quizzes that haven't ended yet
		if now.Before(q.EndTime) {
			availableQuizzes = append(availableQuizzes, q)
		}
	}
	
	return availableQuizzes, nil
}

func (s *quizService) GetQuizForTeacher(ctx context.Context, id uint) (*domain.Quiz, error) {
	return s.quizRepo.GetByID(ctx, id)
}
