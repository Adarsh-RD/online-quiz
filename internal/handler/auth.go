package handler

import (
	"context"
	"net/http"
	"online-quiz/internal/domain"
	"online-quiz/internal/response"
	"online-quiz/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string      `json:"username" binding:"required"`
		Password string      `json:"password" binding:"required"`
		Role     domain.Role `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to register user")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"message": "user created successfully", "user_id": user.ID})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"token": token})
}
