package service

import (
	"errors"
	"math/rand"
	"time"

	"online-quiz/internal/domain"
	"online-quiz/internal/repository"
)

type SessionService interface {
	StartSession(ctx context.Context, quizID, studentID uint) (*domain.QuizSession, error)
	SubmitAnswer(ctx context.Context, sessionID, questionID, optionID uint) error
	SubmitQuiz(ctx context.Context, sessionID uint) error
	HandleTabSwitch(ctx context.Context, sessionID uint) error
	GetSession(ctx context.Context, sessionID uint) (*domain.QuizSession, error)
}

type sessionService struct {
	sessionRepo repository.SessionRepository
	quizRepo    repository.QuizRepository
}

func NewSessionService(sessionRepo repository.SessionRepository, quizRepo repository.QuizRepository) SessionService {
	return &sessionService{sessionRepo, quizRepo}
}

func (s *sessionService) StartSession(ctx context.Context, quizID, studentID uint) (*domain.QuizSession, error) {
	// Check if quiz exists and is available
	quiz, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		return nil, errors.New("quiz not found")
	}

	now := time.Now()
	if now.Before(quiz.StartTime) || now.After(quiz.EndTime) {
		return nil, errors.New("quiz is not currently active")
	}

	// Check if session already exists
	existingSession, err := s.sessionRepo.GetSessionByQuizAndStudent(ctx, quizID, studentID)
	if err == nil && existingSession != nil {
		return nil, errors.New("student already has a session for this quiz")
	}

	session := &domain.QuizSession{
		QuizID:    quizID,
		StudentID: studentID,
		StartTime: &now,
		State:     domain.StateActive,
	}

	// Shuffle questions
	questions := make([]domain.Question, len(quiz.Questions))
	copy(questions, quiz.Questions)
	shuffleQuestions(questions)

	for qIdx, q := range questions {
		sq := domain.SessionQuestion{
			QuestionID: q.ID,
			OrderIndex: qIdx,
		}

		options := make([]domain.Option, len(q.Options))
		copy(options, q.Options)
		shuffleOptions(options)

		for oIdx, o := range options {
			sq.SessionOptions = append(sq.SessionOptions, domain.SessionOption{
				OptionID:   o.ID,
				OrderIndex: oIdx,
			})
		}
		session.SessionQuestions = append(session.SessionQuestions, sq)
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionService) GetSession(ctx context.Context, sessionID uint) (*domain.QuizSession, error) {
	return s.sessionRepo.GetSessionByID(ctx, sessionID)
}

func (s *sessionService) SubmitAnswer(ctx context.Context, sessionID, questionID, optionID uint) error {
	sess, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return errors.New("session not found")
	}

	if sess.State != domain.StateActive {
		return errors.New("cannot submit answer, session is not active")
	}

	// Validate question belongs to session
	var sessionQuestionID uint
	found := false
	for _, sq := range sess.SessionQuestions {
		if sq.QuestionID == questionID {
			sessionQuestionID = sq.ID
			found = true
			break
		}
	}
	if !found {
		return errors.New("question not part of this quiz session")
	}

	ans := domain.SessionAnswer{
		SessionQuestionID: sessionQuestionID,
		SelectedOptionID:  optionID,
	}

	return s.sessionRepo.CreateSessionAnswers(ctx, []domain.SessionAnswer{ans})
}

func (s *sessionService) SubmitQuiz(ctx context.Context, sessionID uint) error {
	return s.sessionRepo.WithTransaction(ctx, func(txRepo repository.SessionRepository) error {
		sess, err := txRepo.GetSessionByID(ctx, sessionID)
		if err != nil {
			return err
		}

		if sess.State != domain.StateActive {
			return errors.New("session is not active")
		}

		// Pull full quiz to score
		quiz, err := s.quizRepo.GetByID(ctx, sess.QuizID)
		if err != nil {
			return err
		}

		score := float64(0)
		// Basic O(N^2) scoring logic, or maps for O(N)
		correctOptions := make(map[uint]float64)
		for _, q := range quiz.Questions {
			for _, o := range q.Options {
				if o.IsCorrect {
					correctOptions[o.ID] = q.Marks
				}
			}
		}

		for _, sq := range sess.SessionQuestions {
			if sq.Answer != nil {
				if marks, ok := correctOptions[sq.Answer.SelectedOptionID]; ok {
					score += marks
				}
			}
		}

		// Apply anti-cheat deduction (1 mark per tab switch)
		penalty := float64(sess.TabSwitchCount)
		finalScore := score - penalty
		if finalScore < 0 {
			finalScore = 0
		}

		sess.Score = finalScore
		sess.State = domain.StateSubmitted
		
		if sess.SuspiciousScore >= 5 { // Threshold hardcoded to 5
			sess.State = domain.StateUnderReview
		}

		return txRepo.UpdateSession(ctx, sess)
	})
}

func (s *sessionService) HandleTabSwitch(ctx context.Context, sessionID uint) error {
	sess, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if sess.State != domain.StateActive {
		return errors.New("session not active")
	}

	sess.TabSwitchCount++
	sess.SuspiciousScore++ // Simplistic 1:1 ratio for now

	return s.sessionRepo.UpdateSession(ctx, sess)
}

func shuffleQuestions(slice []domain.Question) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func shuffleOptions(slice []domain.Option) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
