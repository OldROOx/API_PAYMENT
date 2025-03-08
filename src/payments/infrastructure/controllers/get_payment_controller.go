package controllers

import (
	"app.payment/src/payments/application"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GetPaymentController struct {
	getPaymentUseCase *application.GetPaymentUseCase
}

func NewGetPaymentController(getPaymentUseCase *application.GetPaymentUseCase) *GetPaymentController {
	return &GetPaymentController{getPaymentUseCase: getPaymentUseCase}
}

func (c *GetPaymentController) HandleGetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	payment, err := c.getPaymentUseCase.ExecuteByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}

func (c *GetPaymentController) HandleGetByOrderID(ctx *gin.Context) {
	orderIDStr := ctx.Param("orderID")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	payment, err := c.getPaymentUseCase.ExecuteByOrderID(uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found for the order"})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
