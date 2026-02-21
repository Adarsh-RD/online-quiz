package handler

import (
	"net/http"
	"online-quiz/internal/response"
	"online-quiz/internal/service"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(sessionService service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService}
}

type StartSessionReq struct {
	QuizID uint `json:"quiz_id" binding:"required"`
}

func (h *SessionHandler) StartSession(c *gin.Context) {
	var req StartSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	studentID, _ := c.Get("userID")

	session, err := h.sessionService.StartSession(c.Request.Context(), req.QuizID, studentID.(uint))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, session)
}

type SubmitAnswerReq struct {
	SessionID  uint `json:"session_id" binding:"required"`
	QuestionID uint `json:"question_id" binding:"required"`
	OptionID   uint `json:"option_id" binding:"required"`
}

func (h *SessionHandler) SubmitAnswer(c *gin.Context) {
	var req SubmitAnswerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.sessionService.SubmitAnswer(c.Request.Context(), req.SessionID, req.QuestionID, req.OptionID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "answer submitted"})
}

type SubmitQuizReq struct {
	SessionID uint `json:"session_id" binding:"required"`
}

func (h *SessionHandler) SubmitQuiz(c *gin.Context) {
	var req SubmitQuizReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.sessionService.SubmitQuiz(c.Request.Context(), req.SessionID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "quiz submitted successfully"})
}

func (h *SessionHandler) TabSwitchEvent(c *gin.Context) {
	var req SubmitQuizReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.sessionService.HandleTabSwitch(c.Request.Context(), req.SessionID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "tab switch recorded"})
}
