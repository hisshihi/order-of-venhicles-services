package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
)

type createSubscriptionRequest struct {
	SelectSubscription string `json:"select_subscription" binding:"required,oneof=14days month year"`
	PromoCode          string `json:"promo_code"`
}

func (server *Server) createSubscription(ctx *gin.Context) {
	var req createSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if user.Role.Role != sqlc.RoleProvider {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("только услугодатели могут оформить подписку")))
		return
	}

	// Проверяем, нет ли уже активной подписки
	activeSubscription, err := server.store.GetActiveSubscriptionForProvider(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if err == nil && activeSubscription.Status.StatusSubscription == sqlc.StatusSubscriptionActive {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("у вас уже есть активная подписка")))
		return
	}

	startDate := time.Now()
	var subscriptionType string
	var endDate time.Time
	var standardPrice float64
	var finalPrice float64

	switch req.SelectSubscription {
	case "14days":
		endDate = startDate.AddDate(0, 0, 14)
		subscriptionType = "14days"
		standardPrice = 5000.0
	case "month":
		endDate = startDate.AddDate(0, 1, 0)
		subscriptionType = "month"
		standardPrice = 10000.0
	case "year":
		endDate = startDate.AddDate(1, 0, 0)
		subscriptionType = "year"
		standardPrice = 100000.0
	}

	finalPrice = standardPrice
	var promoCodeID sql.NullInt64
	var discountPercentage int

	// Гарантированно инициализируем строковые представления цен
	finalPriceString := fmt.Sprintf("%.2f", finalPrice) // Используем .2f для сохранения двух знаков после запятой
	standardPriceString := fmt.Sprintf("%.2f", standardPrice)

	// Обработка промокода, если он есть
	var promoCode sqlc.PromoCode
	if req.PromoCode != "" {
		promoCode, err = server.store.GetPromoCodeByCode(ctx, req.PromoCode)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("указанный промокод не существует")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if time.Now().After(promoCode.ValidUntil) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("действие промокода истекло")))
			return
		}

		if promoCode.CurrentUsages.Valid && promoCode.MaxUsages.Valid {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("промокод уже использован максимальное количество раз")))
			return
		}

		discountPercentage = int(promoCode.DiscountPercentage)
		finalPrice = standardPrice * (1 - float64(discountPercentage)/100)
		finalPriceString = fmt.Sprintf("%.2f", finalPrice) // Обновляем строковое представление
		promoCodeID = sql.NullInt64{Int64: promoCode.ID, Valid: true}
		err := server.updateCurrentUsagePromoCode(ctx, promoCode.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("ошибка при применении промокода, попробуйте позже")))
			return
		}
	} else {
		promoCodeID = sql.NullInt64{Valid: false}
	}

	arg := sqlc.CreateSubscriptionParams{
		ProviderID:       user.ID,
		StartDate:        startDate,
		EndDate:          endDate,
		Status:           sqlc.NullStatusSubscription{StatusSubscription: sqlc.StatusSubscriptionActive, Valid: true},
		SubscriptionType: sql.NullString{String: subscriptionType, Valid: true},
		PromoCodeID:      promoCodeID,
		// Используем правильные форматированные строки для цен
		Price:         sql.NullString{String: finalPriceString, Valid: true},
		OriginalPrice: sql.NullString{String: standardPriceString, Valid: true},
	}

	subscription, err := server.store.CreateSubscription(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Формируем ответ с информацией о скидке
	rsp := gin.H{
		"id":                subscription.ID,
		"provider_id":       subscription.ProviderID,
		"start_date":        subscription.StartDate,
		"end_date":          subscription.EndDate,
		"subscription_type": subscriptionType,
		"status":            subscription.Status,
		"original_price":    standardPrice,
		"final_price":       finalPrice,
		"created_at":        subscription.CreatedAt,
	}

	if discountPercentage > 0 {
		rsp["discount_percentage"] = discountPercentage
		rsp["discount_amount"] = standardPrice - finalPrice
		rsp["promo_code_id"] = promoCode.ID

	}

	ctx.JSON(http.StatusOK, rsp)
}

// Проверка активной подписки
func (server *Server) checkSubscriptionActive(ctx *gin.Context) {
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if user.Role.Role == sqlc.RoleClient {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("только услугодатели могут иметь подписку")))
	}

	subscription, err := server.store.GetActiveSubscriptionForProvider(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"has_active_subscription": false,
				"message":                 "у вас нет активной подписки",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Вычисляем оставшееся время подписки
	remainingDays := int(time.Until(subscription.EndDate).Hours() / 24)

	ctx.JSON(http.StatusOK, gin.H{
		"has_active_subscription": true,
		"subscription_id":         subscription.ID,
		"start_date":              subscription.StartDate,
		"end_date":                subscription.EndDate,
		"remaining_days":          remainingDays,
		"status":                  subscription.Status.StatusSubscription,
		"updated_at":              subscription.UpdatedAt,
	})
}

type subscriptionUpdateResponse struct {
	ID         int64                       `json:"id"`
	ProviderID int64                       `json:"provider_id"`
	StartDate  time.Time                   `json:"start_date"`
	EndDate    time.Time                   `json:"end_date"`
	Status     sqlc.NullStatusSubscription `json:"status"`
	CreatedAt  time.Time                   `json:"created_at"`
	UpdatedAt  time.Time                   `json:"updated_at"`
}

// Обновление подписки
func (server *Server) updateSubscription(ctx *gin.Context) {
	var req createSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if user.Role.Role != sqlc.RoleProvider {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("только услугодатели могут оформить подписку")))
		return
	}

	// Проверяем, нет ли уже активной подписки
	activeSubscription, err := server.store.GetActiveSubscriptionForProvider(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if err == nil && activeSubscription.Status.StatusSubscription == sqlc.StatusSubscriptionActive {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("у вас уже есть активная подписка")))
		return
	}

	startDate := time.Now()
	var endDate time.Time
	var standardPrice float64
	var finalPrice float64
	var subscriptionType string

	switch req.SelectSubscription {
	case "14days":
		endDate = startDate.AddDate(0, 0, 14)
		standardPrice = 5000.0
		subscriptionType = "14days"
	case "month":
		endDate = startDate.AddDate(0, 1, 0)
		standardPrice = 10000.0
		subscriptionType = "month"
	case "year":
		endDate = startDate.AddDate(1, 0, 0)
		standardPrice = 100000.0
		subscriptionType = "year"
	}

	finalPrice = standardPrice
	var promoCodeID sql.NullInt64
	var discountPercentage int

	// Обработка промокода, если он есть
	var promoCode sqlc.PromoCode
	if req.PromoCode != "" {
		promoCode, err = server.store.GetPromoCodeByCode(ctx, req.PromoCode)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("указанный промокод не существует")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if time.Now().After(promoCode.ValidUntil) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("действие промокода истекло")))
			return
		}

		if promoCode.CurrentUsages.Valid && promoCode.MaxUsages.Valid {
			if promoCode.CurrentUsages.Int32 >= promoCode.MaxUsages.Int32 {
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("промокод уже использован максимальное количество раз")))
				return
			}
		}

		discountPercentage = int(promoCode.DiscountPercentage)
		finalPrice = standardPrice * (1 - float64(discountPercentage)/100)
		promoCodeID = sql.NullInt64{Int64: promoCode.ID, Valid: true}

		// Обновляем количество использований промокода
		err := server.updateCurrentUsagePromoCode(ctx, promoCode.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("ошибка при применении промокода, попробуйте позже")))
			return
		}
	} else {
		promoCodeID = sql.NullInt64{Valid: false}
	}

	findSubscription, err := server.store.GetSubscriptionByProviderID(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.UpdateSubscriptionParams{
		ID:            findSubscription.ID,
		ProviderID:    user.ID,
		StartDate:     startDate,
		EndDate:       endDate,
		Status:        sqlc.NullStatusSubscription{StatusSubscription: sqlc.StatusSubscriptionActive, Valid: true},
		SubscriptionType: sql.NullString{String: subscriptionType, Valid: true},
		PromoCodeID:   promoCodeID,
		Price:         sql.NullString{String: fmt.Sprintf("%.2f", finalPrice), Valid: true},
		OriginalPrice: sql.NullString{String: fmt.Sprintf("%.2f", standardPrice), Valid: true},
	}

	subscription, err := server.store.UpdateSubscription(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := subscriptionUpdateResponse{
		ID:         subscription.ID,
		ProviderID: subscription.ProviderID,
		StartDate:  subscription.StartDate,
		EndDate:    subscription.EndDate,
		Status:     subscription.Status,
		CreatedAt:  subscription.CreatedAt,
		UpdatedAt:  subscription.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// Просмотр списка подписок

type listSubsciptionsProviderIDRequest struct {
	ProviderID int64 `form:"provider_id" binding:"min=1,required"`
	PageID     int32 `form:"page_id" binding:"min=1,required"`
	PageSize   int32 `form:"page_size" binding:"min=5,max=10,required"`
}

func (server *Server) listSubsciptionsByProviderID(ctx *gin.Context) {
	var req listSubsciptionsProviderIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.ListSubscriptionsByProviderIDParams{
		ProviderID: req.ProviderID,
		Limit:      int64(req.PageSize),
		Offset:     int64((req.PageID - 1) * req.PageSize),
	}

	subscriptions, err := server.store.ListSubscriptionsByProviderID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subscriptions)
}

// Запуск переодической проверки истекших подписок
func (server *Server) startSubscriptionChecker() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			ctx := &gin.Context{}
			err := server.updateExpiredSubscriptions(ctx)
			if err != nil {
				log.Printf("Ошибка при обновлении подписок %v", err)
			}
		}
	}()

	log.Println("Запущена проверка истекших подписок")
}

// Middleware для проверки активной подписки
func (server *Server) subscriptionCheckMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := server.getUserDataFromToken(ctx)
		if err != nil {
			ctx.Abort()
			return
		}

		// Проверка только для провайдеров
		if user.Role.Role != sqlc.RoleProvider {
			ctx.Next()
			return
		}

		// Получаем активную подписку
		subscription, err := server.store.GetActiveSubscriptionForProvider(ctx, user.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusPaymentRequired, gin.H{
					"error": "у вас нет активной подписки, оформите подписку для доступа к этому функционалу",
				})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			ctx.Abort()
			return
		}

		// Проверяем статус подписки
		if subscription.Status.StatusSubscription != sqlc.StatusSubscriptionActive || time.Now().After(subscription.EndDate) {
			ctx.JSON(http.StatusPaymentRequired, gin.H{
				"error": "ваша подписка неактивна, обновите подписку для доступа к этому функционалу",
			})
			ctx.Abort()
			return
		}

		// Добавляем информацию о подписке в контекст
		remainingDays := int(time.Until(subscription.EndDate).Hours() / 24)
		ctx.Set("subscription", gin.H{
			"subscription_id": subscription.ID,
			"end_date":        subscription.EndDate,
			"remaining_days":  remainingDays,
		})

		ctx.Next()
	}
}
