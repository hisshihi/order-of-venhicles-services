package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/hisshihi/order-of-venhicles-services/internal/db"
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

func (server *Server) deleteSubcategoryHandler(ctx *gin.Context) {
    // Получаем ID подкатегории из параметров URL
    subcategoryIDStr := ctx.Param("id")
    subcategoryID, err := strconv.Atoi(subcategoryIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID подкатегории"})
        return
    }

    // Параметры для транзакции
    arg := db.DeleteSubcategoryTxParams{
        SubcategoryID: int64(subcategoryID),
    }

    // Выполняем транзакцию
    result, err := server.store.DeleteSubcategoryTx(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if !result.DeletedSubcategory {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Подкатегория не найдена"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Подкатегория успешно удалена",
        "deleted_data": gin.H{
            "orders": result.DeletedOrders,
            "services": result.DeletedServices,
        },
    })
}


