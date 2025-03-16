package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

type createServiceRequest struct {
	CategoryID  int64  `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description" binding:"required,min=50"`
	Price       string `json:"price" binding:"required"`
}

type createServiceResponse struct {
	ProviderID  int64  `json:"provider_id"`
	CategoryID  int64  `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

func (server *Server) createService(ctx *gin.Context) {
	var req createServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	arg := sqlc.CreateServiceParams{
		ProviderID:  user.ID,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
	}

	service, err := server.store.CreateService(ctx, arg)
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

	rsp := createServiceResponse{
		ProviderID:  service.ProviderID,
		CategoryID:  service.CategoryID,
		Title:       service.Title,
		Description: service.Description,
		Price:       service.Price,
	}

	ctx.JSON(http.StatusOK, rsp)
}
