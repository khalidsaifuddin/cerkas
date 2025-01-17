// fake commit
package middleware

import (
	"strings"

	"github.com/cerkas/cerkas-backend/config"
	"github.com/cerkas/cerkas-backend/core/module"
	"github.com/cerkas/cerkas-backend/handler/api"
	"github.com/cerkas/cerkas-backend/pkg/conn"
	catalogrepository "github.com/cerkas/cerkas-backend/repository/catalog_repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(cfg config.Config, db *gorm.DB) (*gin.Engine, conn.CacheService) {
	if strings.EqualFold(cfg.Environment, "production") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(CORSMiddleware())

	coreRedis, _ := conn.InitRedis(cfg)

	// repository
	catalogRepo := catalogrepository.New(cfg, db)

	// usecase
	catalogUc := module.NewCatalogUsecase(cfg, catalogRepo)

	// handler
	httpHandler := api.NewHTTPHandler(cfg, catalogUc)

	router.POST("t/:tenant_code/p/:product_code/o/:object_code/data", httpHandler.GetObjectData)
	router.POST("t/:tenant_code/p/:product_code/o/:object_code/data/raw", httpHandler.GetDataByRawQuery)
	router.POST("t/:tenant_code/p/:product_code/o/:object_code/data/detail/:serial", httpHandler.GetObjectDetail)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "Page not found"})
	})

	return router, coreRedis
}
