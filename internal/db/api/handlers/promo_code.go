package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

func (server *Server) createPromoCode(ctx *gin.Context) {
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if user.Role.Role != sqlc.RolePartner {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("только партнёр может создавать промокоды")))
		return
	}

	startDate := time.Now()
	id := uuid.New()

	arg := sqlc.CreatePromoCodeParams{
		PartnerID:          user.ID,
		Code:               id.String()[:8],
		DiscountPercentage: 14,
		ValidUntil:         startDate.AddDate(0, 0, 2),
		MaxUsages:          sql.NullInt32{Int32: 1, Valid: true},
		CurrentUsages:      sql.NullInt32{Int32: 0, Valid: false},
	}

	promoCode, err := server.store.CreatePromoCode(ctx, arg)
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

	ctx.JSON(http.StatusOK, promoCode)
}

// Обновление кол-ва использований промокода
func (server *Server) updateCurrentUsagePromoCode(ctx *gin.Context, promoCodeID int64) error {

	arg := sqlc.UpdatePromoCodeByIDParams{
		ID:            promoCodeID,
		CurrentUsages: sql.NullInt32{Int32: 1, Valid: true},
	}

	err := server.store.UpdatePromoCodeByID(ctx, arg)
	if err != nil {
		errorResponse(errors.New("ошибка при обновлении кол-ва использований промокода"))
		return err
	}

	return nil
}

// TODO: реализовать получение всех провайдеров использующих промокод партнёра
