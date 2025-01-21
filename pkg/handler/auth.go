package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polyk005/notesync"
)

func (h *Handler) signUp(c *gin.Context) {
	var input notesync.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

type forgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) forgotPassword(c *gin.Context) {
	var input forgotPasswordInput

	// Проверка на ошибки при связывании JSON
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Здесь предполагается, что у вас есть доступ к h.services.Authorization
	// и SendPasswordResetEmail принимает только email
	err := h.services.Authorization.SendPasswordResetEmail(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}
