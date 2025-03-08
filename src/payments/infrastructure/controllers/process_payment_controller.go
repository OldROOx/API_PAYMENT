package controllers

import (
	"app.payment/src/payments/application"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProcessPaymentRequest struct {
	OrderID uint    `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
	Method  string  `json:"method" binding:"required"`
}

type ProcessPaymentController struct {
	processPaymentUseCase *application.ProcessPaymentUseCase
}

func NewProcessPaymentController(processPaymentUseCase *application.ProcessPaymentUseCase) *ProcessPaymentController {
	return &ProcessPaymentController{processPaymentUseCase: processPaymentUseCase}
}

func (c *ProcessPaymentController) Handle(ctx *gin.Context) {
	var req ProcessPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := c.processPaymentUseCase.Execute(req.OrderID, req.Amount, req.Method)
	if err != nil {
		// El estado depende de si el pago falló por validación o por el procesador
		statusCode := http.StatusInternalServerError
		if payment != nil && payment.Status == "FAILED" {
			statusCode = http.StatusPaymentRequired
		}

		ctx.JSON(statusCode, gin.H{
			"error":   err.Error(),
			"payment": payment,
		})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
