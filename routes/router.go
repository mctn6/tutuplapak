package routes

import (
	"database/sql"
	"tutuplapak/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	router := gin.Default()

	return router
}
