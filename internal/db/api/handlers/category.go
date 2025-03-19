package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

type categoryRequest struct {
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon" binding:"required"`
	Description string `json:"description" binding:"required"`
	Slug string `json:"slug" binding:"required,alphanum"`
}

type categoryRespons struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Description string `json:"description"`
	Slug string `json:"slug"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req categoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.CreateServiceCategoryParams{
		Name: req.Name,
		Icon: req.Icon,
		Description: req.Description,
		Slug: req.Slug,
	}

	category, err := server.store.CreateServiceCategory(ctx, arg)
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

	rsp := categoryRespons{
		ID: category.ID,
		Name: category.Name,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getCategoryByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCategoryByID(ctx *gin.Context) {
	var req getCategoryByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetServiceCategoryByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type listCategoryRequest struct {
	PageID int32 `form:"page_id" binding:"min=1,required"`
	PageSize int32 `form:"page_size" binding:"min=5,max=10,required"`
}

func (server *Server) listCategory(ctx *gin.Context) {
	categories, err := server.store.ListServiceCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

type updateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon" binding:"required"`
	Description string `json:"description" binding:"required"`
	Slug string `json:"slug" binding:"required,alphanum"`
}

func (server *Server) updateCategory(ctx *gin.Context) {
	var reqID getCategoryByIDRequest
	var req updateCategoryRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetServiceCategoryByID(ctx, reqID.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.UpdateServiceCategoryParams{
		ID: category.ID,
		Name: req.Name,
		Icon: req.Icon,
		Description: req.Description,
		Slug: req.Slug,
	}

	category, err = server.store.UpdateServiceCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	var req getCategoryByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteServiceCategory(ctx, req.ID)
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