package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

type createSupportMessageRequest struct {
	SenderID int64 `json:"sender_id" binding:"required,min=1"`
	Subject string `json:"subject" binding:"required"`
	Messages string `json:"messages" binding:"required"`
}

func (server *Server) createSupportMessage(ctx *gin.Context) {
	var req createSupportMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.CreateSupportMessageParams{
		SenderID: req.SenderID,
		Subject: req.Subject,
		Messages: req.Messages,
	}

	supportMessage, err := server.store.CreateSupportMessage(ctx, arg)
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

	ctx.JSON(http.StatusOK, supportMessage)
}

type listSupportMessagesRequest struct {
	PageSize   int32 `form:"page_size" binding:"min=5,max=10,required"`
	PageID     int32 `form:"page_id" binding:"min=1,required"`
}

func (server *Server) listSupportMessages(ctx *gin.Context) {
	var req listSupportMessagesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.ListSupportMessagesParams{
		Limit: int64(req.PageSize),
		Offset: int64((req.PageID - 1) * req.PageSize),
	}

	supportMessages, err := server.store.ListSupportMessages(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	supportMessagesCount, err := server.store.CountSupportMessages(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"support_messages": supportMessages,
		"support_messages_count": supportMessagesCount,
	})
}

type deleteSupportMessagesRequest struct {
	ID   int64 `uri:"id" binding:"min=1,required"`
}

func (server *Server) deleteSupportMessage(ctx *gin.Context) {
	var req deleteSupportMessagesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteSupportMessage(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}