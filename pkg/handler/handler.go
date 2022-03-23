package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/samarec1812/wb-orders-test0/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.LoadHTMLGlob("templates/*")
	api := router.Group("/api")
	{
	//	api.GET("/", h.getAllOrders)
		api.POST("/", h.getOrderById)
		api.POST("/send", h.createOrder)
		api.GET("/", h.getHTML)
	}

	return router
}


