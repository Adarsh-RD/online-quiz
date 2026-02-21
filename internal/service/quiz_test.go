package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"online-quiz/internal/domain"
)

// MockQuizRepository to avoid hitting DB during tests
type MockQuizRepository struct {
	quizzes []domain.Quiz
}

func (m *MockQuizRepository) Create(ctx context.Context, quiz *domain.Quiz) error {
	quiz.ID = uint(len(m.quizzes) + 1)
	m.quizzes = append(m.quizzes, *quiz)
	return nil
}

func (m *MockQuizRepository) GetByID(ctx context.Context, id uint) (*domain.Quiz, error) {
	for _, q := range m.quizzes {
		if q.ID == id {
			return &q, nil
		}
	}
	return nil, errors.New("quiz not found")
}

func (m *MockQuizRepository) GetAll(ctx context.Context) ([]domain.Quiz, error) {
	return m.quizzes, nil
}

func (m *MockQuizRepository) GetPublished(ctx context.Context) ([]domain.Quiz, error) {
	var pub []domain.Quiz
	for _, q := range m.quizzes {
		if q.Published {
			pub = append(pub, q)
		}
	}
	return pub, nil
}

func (m *MockQuizRepository) CreateQuestion(ctx context.Context, question *domain.Question) error {
	return nil
}

func TestCreateQuizValidation(t *testing.T) {
	repo := &MockQuizRepository{}
	svc := NewQuizService(repo)

	now := time.Now()

	tests := []struct {
		name    string
		quiz    *domain.Quiz
		wantErr bool
	}{
		{
			name: "Valid Quiz Time",
			quiz: &domain.Quiz{
				Title:     "Test Quiz",
				StartTime: now,
				EndTime:   now.Add(1 * time.Hour),
				TeacherID: 1,
			},
			wantErr: false,
		},
		{
			name: "Invalid Quiz Time (End before Start)",
			quiz: &domain.Quiz{
				Title:     "Test Quiz Invalid",
				StartTime: now.Add(1 * time.Hour),
				EndTime:   now,
				TeacherID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateQuiz(context.Background(), tt.quiz)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateQuiz() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
