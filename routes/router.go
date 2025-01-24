package routes

import (
	"database/sql"
	"tutuplapak/config"
	v1Handlers "tutuplapak/handlers/v1"
	"tutuplapak/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	router := gin.Default()
	jwtMiddleware := middleware.JWTAuth()

	v1Group := router.Group("/v1")

	productHandler := v1Handlers.NewProductHandler(db)

	productRouter := v1Group.Group("product")
	productRouter.Use(jwtMiddleware)
	productRouter.POST("/", productHandler.CreateProduct)
	productRouter.GET("/", productHandler.GetProducts)
	productRouter.PATCH("/:productId", productHandler.UpdateProduct)
	productRouter.DELETE("/:productId", productHandler.DeleteProduct)

	return router
}
