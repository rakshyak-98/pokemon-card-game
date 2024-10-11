package handlers

import (
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e35562c ( user Registration and Authentication using jwt)
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
<<<<<<< HEAD
<<<<<<< HEAD
	user := &models.User{}
=======
	var user models.User
>>>>>>> 7ec283d (fix conflicts)
=======
	user := &models.User{}
>>>>>>> e35562c ( user Registration and Authentication using jwt)
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
	err := h.userService.RegisterUser(user)
=======
	createUser, err := h.userService.CreateUser(user)
>>>>>>> 7ec283d (fix conflicts)
=======
	err := h.userService.RegisterUser(user)
>>>>>>> e35562c ( user Registration and Authentication using jwt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
	c.JSON(http.StatusCreated, user)
=======
	c.JSON(http.StatusCreated, createUser)
>>>>>>> 7ec283d (fix conflicts)
=======
	c.JSON(http.StatusCreated, user)
>>>>>>> e35562c ( user Registration and Authentication using jwt)
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
