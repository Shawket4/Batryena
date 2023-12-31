package Routes

import (
	"BatrynaBackend/Controllers"
	"BatrynaBackend/Middleware"

	"github.com/gin-gonic/gin"
)

func Setup() {
	const ServerPort string = "3006"
	app := gin.Default()
	// Public Registration And Login
	app.GET("/HelloMessage", Controllers.HelloMessage)
	public := app.Group("/api")
	public.POST("/login", Controllers.Login)
	public.POST("/register", Controllers.Register)
	public.GET("/FetchBranchList", Controllers.FetchBranchList)
	// Protected
	authorized := app.Group("/api/protected")
	authorized.Use(Middleware.JwtAuthMiddleware())
	authorized.GET("/user", Controllers.CurrentUser)
	authorized.GET("/FetchBranches", Controllers.FetchBranches)
	// authorized.GET("/FetchInventories", Controllers.FetchInventories)
	authorized.POST("/RegisterBranch", Controllers.RegisterBranch)
	authorized.POST("/UpdateBranch", Controllers.UpdateBranch)
	authorized.POST("/DeleteBranch", Controllers.DeleteBranch)
	authorized.GET("/FetchBranchData", Controllers.FetchBranchData)
	authorized.GET("/FetchTransactions", Controllers.FetchTransactions)
	public.GET("/FetchBranchesHeatData", Controllers.FetchBranchesHeatData)
	authorized.POST("/RegisterTransaction", Controllers.RegisterTransaction)
	authorized.POST("/UpdateTransaction", Controllers.UpdateTransaction)
	authorized.POST("/DeleteTransaction", Controllers.DeleteTransaction)
	authorized.POST("/RegisterEmployee", Controllers.RegisterEmployee)
	authorized.POST("/GetBranchTransactionsExcel", Controllers.GetBranchTransactionsExcel)
	public.POST("/GenerateShiftOTP", Controllers.GenerateShiftOTP)
	authorized.POST("/SwitchShift", Controllers.SwitchShift)
	if err := app.Run(":" + ServerPort); err != nil {
		panic("Couldn't Start Server On Port " + ServerPort)
	}
}
