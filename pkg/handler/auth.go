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

type UpdateResetPassword struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	var input UpdateResetPassword

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.services.Authorization.UpdatePasswordUser(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"Message": "password update",
	})
}

type requestResetPassword struct {
	Email string `json:"email" binding:"required"`
}

func (h *Handler) requestPasswordReset(c *gin.Context) {
	var input requestResetPassword

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.services.SendPassword.CreateResetToken(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Check to your email",
	})
}

func (h *Handler) ResetPasswordHandler(c *gin.Context) {
	if c.GetString("passwordResetDone") != "" {
		// Если пароль уже был сброшен, перенаправляем на главную страницу
		c.Redirect(http.StatusFound, "/main")
		return
	}

	token := c.Query("token")
	if token == "" {
		newErrorResponse(c, http.StatusBadRequest, "Token is required")
		return
	}

	c.HTML(http.StatusOK, "reset_password.html", gin.H{
		"token": token,
	})
}

func (h *Handler) UpdatePasswordHandler(c *gin.Context) {
	var input struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Пробуем привязать JSON-данные к структуре input
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Здесь вы можете проверить токен на валидность и срок действия
	// Если токен валиден, обновите пароль
	err := h.services.Authorization.UpdatePasswordUserToken(input.Token, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Set("passwordResetDone", true)

	// Возвращаем успешный ответ в формате JSON
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Пароль успешно обновлен",
	})
}
