package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

type reviewCreateRequest struct {
	ProviderID int64  `json:"provider_id" binding:"min=1,required"`
	Rating  int32  `json:"rating" binding:"min=1,max=5,required"`
	Comment string `json:"comment" binding:"required,min=10"`
}

type reviewCreateResponse struct {
	ID         int64  `json:"id"`
	ClientID   int64  `json:"client_id"`
	ProviderID int64  `json:"provider_id"`
	Rating     int32  `json:"rating"`
	Comment    string `json:"comment"`
}

func (server *Server) createReview(ctx *gin.Context) {
	var req reviewCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	arg := sqlc.CreateReviewParams{
		ClientID:   user.ID,
		ProviderID: req.ProviderID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	if arg.ClientID == arg.ProviderID {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("нельзя оставлять отзывы на самого себя")))
		return
	}

	review, err := server.store.CreateReview(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden,
					errorResponse(errors.New("вы уже оставили отзыв на этот заказ")))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := reviewCreateResponse{
		ID:         review.ID,
		ClientID:   review.ClientID,
		ProviderID: review.ProviderID,
		Rating:     review.Rating,
		Comment:    review.Comment,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listReviewByProviderIDRequest struct {
	ProviderID int64 `form:"provider_id" binding:"min=1,required"`
	PageSize   int32 `form:"page_size" binding:"min=5,max=10,required"`
	PageID     int32 `form:"page_id" binding:"min=1,required"`
}

func (server *Server) listReviewByProviderID(ctx *gin.Context) {
	var req listReviewByProviderIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.GetReviewsByProviderIDParams{
		ProviderID: req.ProviderID,
		Limit:      int64(req.PageSize),
		Offset:     int64((req.PageID - 1) * req.PageSize),
	}

	reviews, err := server.store.GetReviewsByProviderID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}

// API получения средней оценки об услугодателе
type getAverageRatingForProviderRequest struct {
	ID int64 `uri:"id" binding:"min=1,required"`
}

func (server *Server) getAverageRatingForProvider(ctx *gin.Context) {
	var req getAverageRatingForProviderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.store.GetAverageRatingForProvider(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	value, err := strconv.ParseFloat(result.AverageRating.String, 64)
	if err != nil {
		fmt.Println("Ошибка преобразования", err)
		return
	}

	averageRating := fmt.Sprintf("%.1f", value)

	result.AverageRating.String = averageRating

	ctx.JSON(http.StatusOK, result)
}

// Удаление отзыва (только если его удаляет автор)
type deleteReviewRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteReview(ctx *gin.Context) {
	var req deleteReviewRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.GetReviewByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if review.ClientID == user.ID || user.Role.Role == sqlc.RoleAdmin {
		arg := sqlc.DeleteReviewParams{
			ID:       req.ID,
			ClientID: review.ClientID,
		}
		err = server.store.DeleteReview(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("нельзя удалить чужой отзыв")))
	}
}

// Просмотр отзывов на себя для услугодателя
type getReviewsByOnlyProviderID struct {
	PageSize int32 `form:"page_size" binding:"min=5,max=10,required"`
	PageID   int32 `form:"page_id" binding:"min=1,required"`
}

func (server *Server) getReviewsByThisProviderID(ctx *gin.Context) {
	var req getReviewsByOnlyProviderID
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if user.Role.Role == sqlc.RoleProvider {
		arg := sqlc.GetReviewsByProviderIDParams{
			ProviderID: user.ID,
			Limit:      int64(req.PageSize),
			Offset:     int64((req.PageID - 1) * req.PageSize),
		}

		reviews, err := server.store.GetReviewsByProviderID(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, reviews)
	} else {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("только услугодатель может смотреть свои отзывы")))
	}
}
