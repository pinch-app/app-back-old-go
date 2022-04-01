package controller

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
)

// Build Gin Engile with CORS
func BuildGinEngine() *gin.Engine {
	// Build Gin Engine
	router := gin.Default()

	// Set Gin Log Mode
	if domain.Config().Server.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Enable CORS with Default Configurations
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
	}
	config.AllowOrigins = []string{
		domain.Config().Url.FrontEndUrl,
		domain.Config().Url.BackEndUrl,
		domain.Config().Url.AdminUrl,
	}
	router.Use(cors.New(config))
	ginMid := &Gin{}
	router.Use(ginMid.DecodeToken)
	return router
}

type Gin struct {
	user interfaces.UserRepo
}

func (Gin) DecodeToken(ctx *gin.Context) {
	id := uint64(2)
	ctx.Set("user", &models.User{
		ID: &id,
	})
}

func getPinchUserFromContext(ctx *gin.Context) (*models.User, error) {
	userVal, ok := ctx.Get("user")
	if !ok {
		return nil, errors.New("user not found in context")
	}
	user, ok := userVal.(*models.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
