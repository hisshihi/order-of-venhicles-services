package api

import "github.com/gin-gonic/gin"

type createOrderResponseRequest struct {
	OrderID int64 `json:"order_id" binding:"required"`
	Message string `json:"message" binding:"required"`
	OfferedPrice string `json:"offered_price" binding:"required"`

}

func (server *Server) createOrderResponse(ctx *gin.Context) {

}