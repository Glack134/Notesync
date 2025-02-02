package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polyk005/notesync/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Загрузка HTML-шаблонов
	router.LoadHTMLGlob("templates/*")

	// Обслуживание статических файлов
	router.Static("/static", "./static")

	router.GET("/main", h.mainHandler)

	//остальные маршруты для блокнота
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-UpdatePassword", h.UpdatePassword)
		auth.POST("/request-password-reset", h.requestPasswordReset)
	}

	help := router.Group("/help")
	{
		help.GET("/reset-password", h.ResetPasswordHandler)
		help.POST("/reset-password", h.UpdatePasswordHandler)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}
		items := api.Group("items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}

func (h *Handler) mainHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
