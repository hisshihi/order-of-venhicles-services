package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
)

type createSubscriptionRequest struct {
	SelectSubscription string `json:"select_subscription" binding:"required,oneof=14days month year"`
}

type subscriptionResponse struct {
	ID         int64                       `json:"id"`
	ProviderID int64                       `json:"provider_id"`
	StartDate  time.Time                   `json:"start_date"`
	EndDate    time.Time                   `json:"end_date"`
	Status     sqlc.NullStatusSubscription `json:"status"`
	CreatedAt  time.Time                   `json:"created_at"`
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
	var endDate time.Time

	switch req.SelectSubscription {
	case "14days":
		endDate = startDate.AddDate(0, 0, 14)
	case "month":
		endDate = startDate.AddDate(0, 1, 0)
	case "year":
		endDate = startDate.AddDate(1, 0, 0)
	}

	arg := sqlc.CreateSubscriptionParams{
		ProviderID: user.ID,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     sqlc.NullStatusSubscription{StatusSubscription: sqlc.StatusSubscriptionActive, Valid: true},
	}

	subscription, err := server.store.CreateSubscription(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := subscriptionResponse{
		ID:         subscription.ID,
		ProviderID: subscription.ProviderID,
		StartDate:  subscription.StartDate,
		EndDate:    subscription.EndDate,
		Status:     subscription.Status,
		CreatedAt:  subscription.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// Проверка активной подписки
func (server *Server) checkSubscriptionActive(ctx *gin.Context) {
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if user.Role.Role != sqlc.RoleProvider {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("только услугодатели могут иметь подписку")))
	}

	subscription, err := server.store.GetActiveSubscriptionForProvider(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
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
