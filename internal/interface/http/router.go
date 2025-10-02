package server

import (
	payment "go-blocker/internal/interface/http/handler"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

// @BasePath /api/v1
func RegisterRoutes(h *payment.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/health", h.HealthHandler)
			pay := v1.Group("/payment")
			{
				// pay.POST("/webhook", h.Webhook)
				// pay.GET("/status/:id", h.Status)
				pay.POST("/check/transaction", h.CheckTx)
				pay.POST("/find/transaction", h.FindLatestTx)
			}
		}
	}
	return r
}
