package handler

import (
	"strconv"
	"weather/api/model"
	"weather/config"
	"weather/pkg/logger"
	"weather/service"
	"weather/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Store    storage.IStorage
	Services service.IServiceMangaer
	Log      logger.ILogger
}

func NewStrg(services service.IServiceMangaer, log logger.ILogger) Handler {
	return Handler{
		Services: services,
		Log: log,
	}
}

func handleResponseLog(c *gin.Context, log logger.ILogger, msg string, statusCode int, data interface{}) {
	resp := model.Response{}

	if statusCode >= 100 && statusCode <= 199 {
		resp.Description = config.ERR_INFORMATION
	} else if statusCode >= 200 && statusCode <= 299 {
		resp.Description = config.SUCCESS
		log.Info("request suucced ", logger.Any("msg:", msg), logger.Int("status:", statusCode))
	} else if statusCode >= 300 && statusCode <= 399 {
		resp.Description = config.ERR_REDIRECTION
	} else if statusCode >= 400 && statusCode <= 499 {
		resp.Description = config.ERR_BADREQUEST
		log.Error("!!!!!!!! BAD REQUEST !!!!!!!!", logger.Any("error: ", msg), logger.Int("status: ", statusCode))
	} else {
		resp.Description = config.ERR_INTERNAL_SERVER
		log.Error("!!!!!!!! ERR_INTERNAL_SERVER !!!!!!!!", logger.Any("error: ", msg), logger.Int("status: ", statusCode))

	}

	resp.StatusCode = statusCode
	resp.Data = data

	c.JSON(resp.StatusCode, resp)
}

func ParsePageQueryParam(c *gin.Context) (uint64, error) {
	pageStr := c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.ParseUint(pageStr, 10, 30)
	if err != nil {
		return 0, err
	}
	if page == 0 {
		return 1, nil
	}
	return page, nil
}

func ParseLimitQueryParam(c *gin.Context) (uint64, error) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.ParseUint(limitStr, 10, 30)
	if err != nil {
		return 0, err
	}

	if limit == 0 {
		return 10, nil
	}
	return limit, nil
}
