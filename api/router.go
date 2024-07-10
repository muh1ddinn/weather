package api

import (
	"weather/pkg/logger"
	"weather/service"

	"github.com/gin-gonic/gin"

	_ "weather/api/docs"
	"weather/api/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// New ...
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(services service.IServiceMangaer, log logger.ILogger) *gin.Engine {
	h := handler.NewStrg(services, log)
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//r.Use(authMiddleware)

	r.GET("/weather/", h.Getweather)
	r.GET("/weatherget/", h.GetAllWeather)

	return r
}

//func authMiddleware(c *gin.Context) {
//	auth := c.GetHeader("Authorization")
//	if auth == "" {
//		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized, you can use any character for ApiKeyAuth (apiKey)"))
//	}
//	c.Next()
//}
