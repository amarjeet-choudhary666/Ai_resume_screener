package resumeRouter

import (
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/api/controller"
	"github.com/gin-gonic/gin"
)

func ResumeRoutes(r *gin.Engine) {
	resumeGroup := r.Group("/resume")
	{
		resumeGroup.POST("/upload", controller.UploadResume)
	}
}
