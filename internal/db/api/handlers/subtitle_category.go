package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

type createSubtitleCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createSubtitleCategory(ctx *gin.Context) {
	var req createSubtitleCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subtitleCategory, err := server.store.CreateSubtitle(ctx, req.Name)
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

	ctx.JSON(http.StatusOK, gin.H{
		"id": subtitleCategory.ID,
		"name": subtitleCategory.Name,
	})
}

func (server *Server) listSubtitleCategory(ctx *gin.Context) {
	subtitleCategories, err := server.store.ListSubtitleCategory(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subtitleCategories)
}

func (server *Server) updateSubtitleCategory(ctx *gin.Context) {
	var req struct {
		ID int64 `form:"id" binding:"required,min=1"`
		Name string `form:"name" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.UpdateSubtitleCategoryParams{
		ID: req.ID,
		Name: req.Name,
	}

	subtitleCategory, err := server.store.UpdateSubtitleCategory(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": subtitleCategory.ID,
		"name": subtitleCategory.Name,
	})
}

func (server *Server) deleteSubtitleCategory(ctx *gin.Context) {
	var req struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteSubtitleCategory(ctx, req.ID)
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


