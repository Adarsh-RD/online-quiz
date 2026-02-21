package handler

import (
	"net/http"
	"online-quiz/internal/domain"
	"online-quiz/internal/response"
	"online-quiz/internal/service"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService service.QuizService
}

func NewQuizHandler(quizService service.QuizService) *QuizHandler {
	return &QuizHandler{quizService}
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var req domain.Quiz
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	teacherID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	req.TeacherID = teacherID.(uint)

	if err := h.quizService.CreateQuiz(c.Request.Context(), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, req)
}

func (h *QuizHandler) ListQuizzes(c *gin.Context) {
	quizzes, err := h.quizService.GetQuizzesForStudent(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch quizzes")
		return
	}
	response.Success(c, http.StatusOK, quizzes)
}
