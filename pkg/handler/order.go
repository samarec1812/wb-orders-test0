package handler

import (
	"github.com/gin-gonic/gin"
	orders "github.com/samarec1812/wb-orders-test0"
	"net/http"
	"strconv"
)

func (h *Handler) createOrder(c *gin.Context) {
	var input orders.Orders
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Order.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllOrdersResponse struct {
	Data []orders.Orders `json:"data"`
}

func (h *Handler) getOrderById(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		newErrorHTMLResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	order, err := h.services.GetById(id)
	if err != nil {
		newErrorHTMLResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
			"order": order,
	})
}

func (h *Handler) getHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}