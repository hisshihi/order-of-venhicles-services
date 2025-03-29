package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/hisshihi/order-of-venhicles-services/pkg/util"
)

type createPaymentRequest struct {
	Amount             string `json:"amount" binding:"required,min=1"`
	PaymentMethod      string `json:"payment_method" binding:"required"`
	SelectSubscription string `json:"select_subscription" binding:"required,oneof=14days month year contribution"`
	PromoCode          string `json:"promo_code"`
	Payment            string `json:"payment" binding:"required,oneof=buy_sub update_sub"`
}

// TODO: после тестов убрать
// Тестовый эндпоинт для эмулции успешного платежа
func (server *Server) simuldateSuccessfulPayment(ctx *gin.Context) {
	var req struct {
		PaymentID string `json:"payment_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем платёж по ID
	paymentID, err := strconv.ParseInt(req.PaymentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payment, err := server.store.GetPaymentByID(ctx, paymentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Обновляем статус платежа
	updatePaymentArg := sqlc.UpdatePaymentStatusParams{
		ID:     payment.ID,
		Status: sqlc.NullStatusPayment{StatusPayment: sqlc.StatusPaymentCompleted, Valid: true},
	}

	_, err = server.store.UpdatePaymentStatus(ctx, updatePaymentArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Получаем сохранённые детали подписки
	subscriptionDetails, err := server.getSubscriptionDetails(ctx, req.PaymentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Создаём или обновляем подписку
	if subscriptionDetails.IsUpdate {
		err = server.updateSubscriptionAfterPayment(ctx, payment.UserID, subscriptionDetails)
	} else {
		err = server.createSubscriptionAfterPayment(ctx, payment.UserID, subscriptionDetails)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success", 
		"message": "Симуляция успешного платежа выполнена",
	})
}

// Шаг 1: Создание платежа и получение ссылки для оплаты
func (server *Server) initiateSubscriptionPayment(ctx *gin.Context) {
	var req createPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if user.Role.Role == sqlc.RoleClient {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("нет доступа к этой функции")))
		return
	}

	// Расчёт суммы с учётом промокода
	subscriptionDetails, err := server.calculateSubscriptionDetails(ctx, user.ID, req.SelectSubscription, req.PromoCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Создаём уникальный идентификатор платежа
	paymentID := util.GenerateUniqueID()

	// Создаём запись о платеже в статусе "pending"
	arg := sqlc.CreatePaymentParams{
		ID:            paymentID,
		UserID:        user.ID,
		Amount:        fmt.Sprintf("%.2f", subscriptionDetails.FinalPrice),
		PaymentMethod: req.PaymentMethod,
		Status:        sqlc.NullStatusPayment{StatusPayment: sqlc.StatusPaymentPending, Valid: true},
	}

	payment, err := server.store.CreatePayment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Сохраняем данные подписки во временное хранилище
	err = server.storeSubscriptionDetails(ctx, paymentID, user.ID, subscriptionDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// TODO: тут добавить ссылку на оплату
	paymentLink := server.generatePaymentLink(paymentID, subscriptionDetails.FinalPrice)

	ctx.JSON(http.StatusOK, gin.H{
		"payment_id":        payment.ID,
		"amount":            subscriptionDetails.FinalPrice,
		"payment_link":      paymentLink,
		"subscription_type": subscriptionDetails.SubscriptionType,
		"discount_applied":  subscriptionDetails.DiscountApplied,
		"original_price":    subscriptionDetails.OriginalPrice,
		"status":            "pending",
	})

}

// Шаг 2: Обработка результата платежа (вызывается webhook`ом платёжной системы)
func (server *Server) processPaymentCallback(ctx *gin.Context) {
	var req struct {
		PaymentID string `json:"payment_id"`
		Status    string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Проверка подписи запроса для безопасности (должна быть реализована)
	if !server.verifyPaymentCallback(ctx.Request) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("неверная подпись запроса")))
		return
	}

	// Получаем платёж по ID
	paymentID, err := strconv.ParseInt(req.PaymentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payment, err := server.store.GetPaymentByID(ctx, paymentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Обновляем статус платежа
	if req.Status == "success" {
		updatePaymentArg := sqlc.UpdatePaymentStatusParams{
			ID:     payment.ID,
			Status: sqlc.NullStatusPayment{StatusPayment: sqlc.StatusPaymentCompleted, Valid: true},
		}

		_, err = server.store.UpdatePaymentStatus(ctx, updatePaymentArg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Получаем сохранённые детали подписки
		subscriptionDetails, err := server.getSubscriptionDetails(ctx, req.PaymentID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Создаём или обновляем подписку
		if subscriptionDetails.IsUpdate {
			err = server.updateSubscriptionAfterPayment(ctx, payment.UserID, subscriptionDetails)
		} else {
			err = server.createSubscriptionAfterPayment(ctx, payment.UserID, subscriptionDetails)
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		// Обвноляем статус платежа на "failed"
		updatePaymentArg := sqlc.UpdatePaymentStatusParams{
			ID:     payment.ID,
			Status: sqlc.NullStatusPayment{StatusPayment: sqlc.StatusPaymentFailed, Valid: true},
		}

		_, err = server.store.UpdatePaymentStatus(ctx, updatePaymentArg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "failed"})
	}
}

// Проверка статуса платежа (для клиентского приложения)
func (server *Server) checkPaymentStatus(ctx *gin.Context) {
	var req struct {
		PaymentID string `uri:"payment_id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	paymentID, err := strconv.ParseInt(req.PaymentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err := server.store.GetPaymentByID(ctx, paymentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	// Проверяем, что пользователь имеет доступ к этому платежу
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if payment.UserID != user.ID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("нет доступа к данному платежу")))
		return
	}

	// Получаем статус подписки, если платеж успешен
	var subscriptionStatus string
	if payment.Status.StatusPayment == sqlc.StatusPaymentCompleted {
		subscription, err := server.store.GetActiveSubscriptionForProvider(ctx, user.ID)
		if err == nil {
			subscriptionStatus = string(subscription.Status.StatusSubscription)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"payment_id":          payment.ID,
		"status":              payment.Status.StatusPayment,
		"amount":              payment.Amount,
		"created_at":          payment.CreatedAt,
		"subscription_status": subscriptionStatus,
	})
}

// -----------------------------
// Интеграция с внешней платёжной системой
// -----------------------------
// Генерация платежной ссылки (пример для интеграции с Yookassa)
func (server *Server) generatePaymentLink(paymentID int64, amount float64) string {
	// Это заглушка, которую нужно заменить фактической интеграцией с платежной системой

	// Пример интеграции с Yookassa
	// client := yookassa.NewClient(server.config.YookassaShopID, server.config.YookassaSecretKey)
	// payment, err := client.CreatePayment(
	//     amount,
	//     "RUB",
	//     fmt.Sprintf("Оплата подписки #%d", paymentID),
	//     fmt.Sprintf("%s/payments/callback", server.config.BaseURL),
	//     map[string]string{"payment_id": fmt.Sprintf("%d", paymentID)},
	// )

	// В реальном коде здесь нужно вернуть URL для оплаты от платежной системы
	return fmt.Sprintf("https://example.com/pay/%d?amount=%.2f", paymentID, amount)
}

// Верификация callback от платежной системы
func (server *Server) verifyPaymentCallback(r *http.Request) bool {
	// Это заглушка, которую нужно заменить фактической проверкой подписи запроса

	// Пример для Yookassa
	// signature := r.Header.Get("Signature")
	// body, _ := ioutil.ReadAll(r.Body)
	// r.Body = ioutil.NopCloser(bytes.NewReader(body)) // Восстанавливаем Body для дальнейшего чтения
	// computedSignature := hmac.New(sha256.New, []byte(server.config.YookassaSecretKey))
	// computedSignature.Write(body)
	// expectedSignature := hex.EncodeToString(computedSignature.Sum(nil))
	// return signature == expectedSignature

	return true
}

// -----------------------------
// ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
// -----------------------------
// Структура для хранения деталей подписки
type SubscriptionDetails struct {
	SubscriptionType string
	StartDate        time.Time
	EndDate          time.Time
	OriginalPrice    float64
	FinalPrice       float64
	DiscountApplied  float64
	PromoCodeID      int64
	IsUpdate         bool
}

// Расчёт деталей подписки с учётом промокода
func (server *Server) calculateSubscriptionDetails(ctx *gin.Context, userID int64, subscriptionType, promoCode string) (SubscriptionDetails, error) {
	var details SubscriptionDetails
	details.SubscriptionType = subscriptionType
	details.StartDate = time.Now()

	// Определение цены и срока подписки
	switch subscriptionType {
	case "14days":
		details.EndDate = details.StartDate.AddDate(0, 0, 14)
		details.OriginalPrice = 5000.0
	case "month":
		details.EndDate = details.StartDate.AddDate(0, 1, 0)
		details.OriginalPrice = 10000.0
	case "year":
		details.EndDate = details.StartDate.AddDate(1, 0, 0)
		details.OriginalPrice = 100000.0
	case "contribution":
		details.EndDate = details.StartDate.AddDate(120, 0, 0)
		details.OriginalPrice = 20000.0
	default:
		return details, errors.New("некорректный тип подписки")
	}

	details.FinalPrice = details.OriginalPrice

	// Проверка наличия существующей подписки
	subscription, err := server.store.GetSubscriptionByProviderID(ctx, userID)
	if err == nil {
		details.IsUpdate = true
	}
	if err != nil && err != sql.ErrNoRows {
		return details, err
	}
	if err == nil && subscription.Status.StatusSubscription == sqlc.StatusSubscriptionActive {
		return details, errors.New("у вас уже есть активная подписка")
	}

	// Обработка промокода, если он указан
	if promoCode != "" {
		promoCodeObj, err := server.store.GetPromoCodeByCode(ctx, promoCode)
		if err != nil {
			if err == sql.ErrNoRows {
				return details, errors.New("указанный промокод не существует")
			}
			return details, errors.New("указанный промокод не существует")
		}

		if time.Now().After(promoCodeObj.ValidUntil) {
			return details, errors.New("срок действия промокода истёк")
		}

		if promoCodeObj.CurrentUsages.Valid && promoCodeObj.MaxUsages.Valid {
			if promoCodeObj.CurrentUsages.Int32 >= promoCodeObj.MaxUsages.Int32 {
				return details, errors.New("промокод уже использован максимальное кол-во раз")
			}
		}

		// Расчёт скидки
		discountPercentage := float64(promoCodeObj.DiscountPercentage)
		details.DiscountApplied = details.OriginalPrice * (discountPercentage / 100)
		details.FinalPrice = details.OriginalPrice - details.DiscountApplied
		details.PromoCodeID = promoCodeObj.ID
	}

	return details, nil
}

// Сохранение деталей подписки во временное хранилище
func (server *Server) storeSubscriptionDetails(ctx *gin.Context, paymentID int64, userID int64, details SubscriptionDetails) error {
	// Здесь можно использовать Redis или таблицу в БД для хранения деталей
	// Временное решение - можно использовать локальное хранилище в памяти

	// Пример структуры для таблицы в БД:
	arg := sqlc.CreatePendingSubscriptionParams{
		PaymentID:        paymentID,
		UserID:           userID,
		SubscriptionType: details.SubscriptionType,
		StartDate:        details.StartDate,
		EndDate:          details.EndDate,
		OriginalPrice:    fmt.Sprintf("%.2f", details.OriginalPrice),
		FinalPrice:       fmt.Sprintf("%.2f", details.FinalPrice),
		IsUpdate:         details.IsUpdate,
	}

	if details.PromoCodeID != 0 {
		arg.PromoCodeID = sql.NullInt64{Int64: details.PromoCodeID, Valid: true}
	}

	_, err := server.store.CreatePendingSubscription(ctx, arg)
	return err
}

// Получение сохранённых деталей
func (server *Server) getSubscriptionDetails(ctx *gin.Context, paymentIDStr string) (SubscriptionDetails, error) {
	var details SubscriptionDetails

	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		return details, err
	}

	pendingSub, err := server.store.GetPendingSubscriptionByPaymentID(ctx, paymentID)
	if err != nil {
		return details, err
	}

	details.SubscriptionType = pendingSub.SubscriptionType
	details.StartDate = pendingSub.StartDate
	details.EndDate = pendingSub.EndDate

	origPrice, _ := strconv.ParseFloat(pendingSub.OriginalPrice, 64)
	details.OriginalPrice = origPrice

	finalPrice, _ := strconv.ParseFloat(pendingSub.FinalPrice, 64)
	details.FinalPrice = finalPrice

	details.IsUpdate = pendingSub.IsUpdate

	if pendingSub.PromoCodeID.Valid {
		details.PromoCodeID = pendingSub.PromoCodeID.Int64
	}

	return details, nil
}

// Создание подписки после успешной оплаты
func (server *Server) createSubscriptionAfterPayment(ctx *gin.Context, userID int64, details SubscriptionDetails) error {
	var promoCodeID sql.NullInt64
	if details.PromoCodeID != 0 {
		promoCodeID = sql.NullInt64{Int64: details.PromoCodeID, Valid: true}

		// Увеличиваем счетчик использования промокода
		err := server.updateCurrentUsagePromoCode(ctx, details.PromoCodeID)
		if err != nil {
			return err
		}
	}

	arg := sqlc.CreateSubscriptionParams{
		ProviderID:       userID,
		StartDate:        details.StartDate,
		EndDate:          details.EndDate,
		Status:           sqlc.NullStatusSubscription{StatusSubscription: sqlc.StatusSubscriptionActive, Valid: true},
		SubscriptionType: sql.NullString{String: details.SubscriptionType, Valid: true},
		PromoCodeID:      promoCodeID,
		Price:            sql.NullString{String: fmt.Sprintf("%.2f", details.FinalPrice), Valid: true},
		OriginalPrice:    sql.NullString{String: fmt.Sprintf("%.2f", details.OriginalPrice), Valid: true},
	}

	_, err := server.store.CreateSubscription(ctx, arg)
	return err
}

// Обновление подписки после успешной оплаты
func (server *Server) updateSubscriptionAfterPayment(ctx *gin.Context, userID int64, details SubscriptionDetails) error {
	var promoCodeID sql.NullInt64
	if details.PromoCodeID != 0 {
		promoCodeID = sql.NullInt64{Int64: details.PromoCodeID, Valid: true}

		// Увеличиваем счетчик использования промокода
		err := server.updateCurrentUsagePromoCode(ctx, details.PromoCodeID)
		if err != nil {
			return err
		}
	}

	subscription, err := server.store.GetSubscriptionByProviderID(ctx, userID)
	if err != nil {
		return err
	}

	arg := sqlc.UpdateSubscriptionParams{
		ID:               subscription.ID,
		ProviderID:       userID,
		StartDate:        details.StartDate,
		EndDate:          details.EndDate,
		Status:           sqlc.NullStatusSubscription{StatusSubscription: sqlc.StatusSubscriptionActive, Valid: true},
		SubscriptionType: sql.NullString{String: details.SubscriptionType, Valid: true},
		PromoCodeID:      promoCodeID,
		Price:            sql.NullString{String: fmt.Sprintf("%.2f", details.FinalPrice), Valid: true},
		OriginalPrice:    sql.NullString{String: fmt.Sprintf("%.2f", details.OriginalPrice), Valid: true},
	}

	_, err = server.store.UpdateSubscription(ctx, arg)
	return err
}
