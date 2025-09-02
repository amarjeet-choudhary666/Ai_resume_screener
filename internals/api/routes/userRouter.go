package userRouter

import (
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/api/controller"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", controller.SignUp())
	}
}
