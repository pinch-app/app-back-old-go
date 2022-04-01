package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
)

type UserController struct {
	user interfaces.UserService
}

// NewUserController creates new group for user endpoint
func NewUserController(router *gin.Engine, user interfaces.UserService) {
	cnt := &UserController{
		user: user,
	}

	group := router.Group("/user")
	{
		group.POST("", cnt.Create)
		group.GET("/:id", cnt.FindByID)
		group.GET("", cnt.FindAll)
		group.PUT("/:id", cnt.UpdateByID)
		group.DELETE("/:id", cnt.DeleteByID)
	}

}

func (c *UserController) Create(ctx *gin.Context) {
	user := &models.User{}

	// Bind user data to user struct
	err := ctx.BindJSON(user)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	// Create User
	err = c.user.Create(user)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"status": "ok",
		"user":   user,
	})
}

func (c *UserController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := ParseUint64(id)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	// Find User
	user, err := c.user.FindOne(&models.User{ID: &userID})
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"status": "ok",
		"user":   user,
	})
}

func (c *UserController) FindAll(ctx *gin.Context) {
	users, err := c.user.FindAll()
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"status": "ok",
		"users":  users,
	})
}

func (c *UserController) UpdateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := ParseUint64(id)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	user := &models.User{}
	// Bind user data to user struct
	err = ctx.BindJSON(user)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	user.ID = &userID

	// Update User
	err = c.user.Update(user)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"status": "ok",
		"user":   user,
	})
}

func (c *UserController) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := ParseUint64(id)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	user := &models.User{
		ID: &userID,
	}

	// Delete User
	err = c.user.Delete(user)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"status": "ok",
		"user":   user,
	})
}
