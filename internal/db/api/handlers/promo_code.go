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

type createPromoCodeRequest struct {
	PartnerID          int64  `json:"partner_id" binding:"required,min=1"`
	DiscountPercentage int32  `json:"discount_percentage" binding:"required,min=1"`
	Code               string `json:"code"`
	Valid              int    `json:"valid"`
}

func (server *Server) createPromoCode(ctx *gin.Context) {
	var req createPromoCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if user.Role.Role != sqlc.RoleAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("только админ может создавать промокоды")))
		return
	}

	startDate := time.Now()

	var id string
	if len(req.Code) > 0 {
		id = req.Code
	} else {
		id = uuid.NewString()
	}

	var valiDay int
	if req.Valid > 0 {
		valiDay = req.Valid
	} else {
		valiDay = 2
	}

	arg := sqlc.CreatePromoCodeParams{
		PartnerID:          req.PartnerID,
		Code:               id,
		DiscountPercentage: req.DiscountPercentage,
		ValidUntil:         startDate.AddDate(0, 0, valiDay),
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

type listPromoCodesRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// Вывод всех промокодов
func (server *Server) listPromoCodes(ctx *gin.Context) {
	var req listPromoCodesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.ListPromoCodesParams{
		Limit: int64(req.PageSize),
		Offset: int64((req.PageID - 1) * req.PageSize),
	}

	promoCodes, err := server.store.ListPromoCodes(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	countPromoCodes, err := server.store.CountPromoCode(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"promo_codes": promoCodes,
		"promo_codes_count": countPromoCodes,
	})
}

func (server *Server) listPromoCodesByPartner(ctx *gin.Context) {
	var req listPromoCodesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.ListPromoCodesByPartnerIDParams{
		PartnerID: user.ID,
		Limit: int64(req.PageSize),
		Offset: int64((req.PageID - 1) * req.PageSize),
	}

	promoCodes, err := server.store.ListPromoCodesByPartnerID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	countPromoCodes, err := server.store.CountPromoCode(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"promo_codes": promoCodes,
		"promo_codes_count": countPromoCodes,
	})
}

// TODO: реализовать получение всех провайдеров использующих промокод партнёра
type getAllProvidersByPartnerPromosRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}
func (server *Server) getAllProvidersByPartnerPromos(ctx *gin.Context) {
	var req getAllProvidersByPartnerPromosRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.GetAllProvidersByPartnerPromosParams{
		PartnerID: user.ID,
		Limit: int64(req.PageSize),
		Offset: int64((req.PageID - 1) * req.PageSize),
	}

	listPromoCodesByProviderID, err := server.store.GetAllProvidersByPartnerPromos(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, listPromoCodesByProviderID)
}