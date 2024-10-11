package handlers

import (
	"net/http"

	"github.com/rakshyak-98/pokemonapi/internal/models"
	"github.com/rakshyak-98/pokemonapi/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", h.CreateUser)
		userGroup.GET("/:id", h.GetUser)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
