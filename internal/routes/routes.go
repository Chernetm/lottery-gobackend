package routes

import (
	"lottery-backend/internal/handlers"
	"lottery-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authH *handlers.AuthHandler,
	lotteryH *handlers.LotteryHandler,
	ticketH *handlers.TicketHandler,
	adminH *handlers.AdminHandler,
	withdrawalH *handlers.WithdrawalHandler,
	paymentH *handlers.PaymentHandler,
) {
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
			auth.POST("/admin/register", authH.AdminRegister)
			auth.POST("/admin/login", authH.AdminLogin)

			// Authenticated routes
			authenticated := auth.Group("/")
			authenticated.Use(middleware.AuthMiddleware())
			{
				authenticated.GET("/profile", authH.GetProfile)
				authenticated.POST("/change-password", authH.ChangePassword)
			}
		}

		// Lottery routes
		lotteries := api.Group("/lotteries")
		{
			lotteries.GET("", lotteryH.GetAllLotteries)
			lotteries.GET("/:id", lotteryH.GetLotteryByID)
			lotteries.GET("/items", lotteryH.GetLotteryItems)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			user := protected.Group("/user")
			{
				user.POST("/withdrawals", withdrawalH.RequestWithdrawal)
				user.GET("/withdrawals", withdrawalH.GetMyWithdrawals)
			}

			// Ticket routes
			tickets := protected.Group("/tickets")
			{
				tickets.GET("/my", ticketH.GetUserTickets)
				tickets.POST("/purchase", ticketH.PurchaseTicket)
				tickets.POST("/:id/reveal", ticketH.RevealTicket)
			}

			// Admin routes
			admin := protected.Group("/admin")
			// admin.Use(middleware.AdminOnly()) // Should implement this next
			{
				admin.GET("/profile", authH.GetAdminProfile)
				admin.GET("/stats", adminH.GetAdminStats)
				admin.POST("/items", adminH.CreateItem)
				admin.DELETE("/items/:id", adminH.DeleteItem)
				admin.GET("/tickets", adminH.GetAllTickets)
				admin.PATCH("/tickets/:id/status", adminH.UpdateTicketStatus)
				admin.PUT("/lotteries/:id", adminH.UpdateLottery)
				admin.PATCH("/lotteries/:id/status", adminH.UpdateLotteryStatus)
				admin.POST("/lotteries/:id/draw", adminH.DrawWinner)
				admin.GET("/withdrawals", withdrawalH.GetAllWithdrawals)
				admin.PATCH("/withdrawals/:id/status", withdrawalH.UpdateWithdrawalStatus)
			}

			protected.POST("/lotteries", lotteryH.CreateLottery)

			// User profile
			// protected.GET("/auth/profile", authH.GetProfile)
		}

		// Payment routes
		payments := api.Group("/payments")
		{
			payments.POST("/webhook", paymentH.HandleWebhook)

			authorizedPayments := payments.Group("/")
			authorizedPayments.Use(middleware.AuthMiddleware())
			{
				authorizedPayments.GET("/verify/:tx_ref", paymentH.VerifyTransaction)
			}
		}
	}
}
