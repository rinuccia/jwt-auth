package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rinuccia/jwt-auth/controllers"
)

func InitRoutes(router *gin.Engine) {

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", controllers.SignUp)
		auth.POST("/sign-in", controllers.SignIn)
	}
}
