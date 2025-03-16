package api

import (
	"database/sql"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

type createOrderRequest struct {
	ServiceID int64 `json:"service_id" binding:"required,min=1"`
}

type createOrderResponse struct {
	ID        int64                 `json:"id"`
	ClientID  int64                 `json:"client_id"`
	ServiceID int64                 `json:"service_id"`
	Status    sqlc.NullStatusOrders `json:"status"`
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	arg := sqlc.CreateOrderParams{
		ClientID:  user.ID,
		ServiceID: req.ServiceID,
		Status:    sqlc.NullStatusOrders{StatusOrders: sqlc.StatusOrdersPending, Valid: true},
	}

	order, err := server.store.CreateOrder(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createOrderResponse{
		ID:        order.ID,
		ClientID:  order.ClientID,
		ServiceID: order.ServiceID,
		Status:    order.Status,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getOrderByIDRequest struct {
	ID int64 `uri:"id" binding:"min=1,required"`
}

func (server *Server) getOrderByID(ctx *gin.Context) {
	var req getOrderByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	order, err := server.store.GetOrderByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if order.ClientID != user.ID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("у вас нет прав на просмотр этого заказа")))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

type listOrderRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listOrders(ctx *gin.Context) {
	var req listOrderRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	arg := sqlc.ListOrdersParams{
		ClientID: user.ID,
		Limit: int64(req.PageID),
		Offset: int64((req.PageID - 1) * req.PageSize),
	}

	orders, err := server.store.ListOrders(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
