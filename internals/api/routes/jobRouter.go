package jobRouter

import (
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/api/controller"
	"github.com/gin-gonic/gin"
)

func JobRoutes(r *gin.Engine) {
	jobGroup := r.Group("/job")
	{
		jobGroup.POST("/create", controller.CreateJob)
		jobGroup.GET("/list", controller.GetJobs)
		jobGroup.POST("/match/:jobId", controller.MatchCandidates)
	}
}
