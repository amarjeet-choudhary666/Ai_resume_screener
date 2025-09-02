package routes

import (
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/api/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.Use(
		middlewares.CORSMiddleware(),
		middlewares.LoggerMiddleware(),
		middlewares.RateLimitMiddleware(),
		middlewares.ErrorHandlerMiddleware(),
	)

	UserRoutes(router)
	ResumeRoutes(router)
	JobRoutes(router)
}
