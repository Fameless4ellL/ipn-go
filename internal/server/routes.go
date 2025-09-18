package server

import (
	"net/http"

	docs "go-blocker/cmd/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-blocker/internal/database"
	"go-blocker/internal/payment"

	"github.com/gin-contrib/cors"
)

// @BasePath /api/v1

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", s.healthHandler)

	pay := v1.Group("/payment")
	paymentRepo := database.NewPaymentRepository(s.db)
	paymentService := payment.NewPaymentService(paymentRepo)
	paymentHandler := &payment.Handler{Service: paymentService}
	// pay.POST("/webhook", paymentHandler.Webhook)
	// pay.GET("/status/:id", paymentHandler.Status)
	pay.POST("/check/transaction", paymentHandler.CheckTx)
	pay.POST("/find/transaction", paymentHandler.FindLatestTx)

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"db":     "connected",
	})
}
