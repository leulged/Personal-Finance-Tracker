package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"personal-finance-tracker/domain/entities"
	usecase "personal-finance-tracker/UseCase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user entities.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Clear any role or verification status sent by client
	// These will be set by the server based on business logic
	user.Role = ""
	user.IsVerified = false

	createdUser, err := h.userUsecase.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}