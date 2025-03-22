package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

// createOrderRequest представляет запрос на создание заказа
type createOrderRequest struct {
	ClientMessage string `json:"client_message" binding:"required"` // Сообщение от клиента
	CategoryID    int64  `json:"category_id" binding:"required"`    // ID категории услуги
	OrderDate     string `json:"order_date"`                        // Планируемая дата заказа (необязательно)
}

// createOrderResponse представляет ответ на создание заказа
type createOrderResponse struct {
	ID            int64                 `json:"id"`
	ClientID      int64                 `json:"client_id"`
	CategoryID    int64                 `json:"category_id"`
	Status        sqlc.NullStatusOrders `json:"status"`
	ClientMessage string                `json:"client_message"`
	CategoryName  string                `json:"category_name,omitempty"`
	OrderDate     string                `json:"order_date,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
}

// createOrder обрабатывает запрос на создание нового заказа
func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем данные пользователя из токена
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	// Проверяем существование категории
	category, err := server.store.GetServiceCategoryByID(ctx, req.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("категория услуги не найдена")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Парсим дату заказа, если она указана
	var orderDate sql.NullTime
	if req.OrderDate != "" {
		parsedDate, err := time.Parse("2006-01-02T15:04:05Z", req.OrderDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("неверный формат даты заказа (используйте ISO 8601)")))
			return
		}
		orderDate = sql.NullTime{Time: parsedDate, Valid: true}
	}

	// Создаем параметры для запроса
	arg := sqlc.CreateOrderParams{
		ClientID:      user.ID,
		CategoryID:    req.CategoryID,
		Status:        sqlc.NullStatusOrders{StatusOrders: sqlc.StatusOrdersPending, Valid: true},
		ClientMessage: sql.NullString{String: req.ClientMessage, Valid: true},
		OrderDate:     orderDate,
	}

	// Создаем заказ в базе данных
	order, err := server.store.CreateOrder(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("указана несуществующая категория")))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Формируем ответ
	rsp := createOrderResponse{
		ID:            order.ID,
		ClientID:      order.ClientID,
		CategoryID:    order.CategoryID,
		Status:        order.Status,
		ClientMessage: order.ClientMessage.String,
		CategoryName:  category.Name,
		CreatedAt:     order.CreatedAt,
	}

	if order.OrderDate.Valid {
		rsp.OrderDate = order.OrderDate.Time.Format("2006-01-02T15:04:05Z")
	}

	ctx.JSON(http.StatusOK, rsp)
}

// getOrderByIDRequest представляет запрос на получение заказа по ID
type getOrderByIDRequest struct {
	ID int64 `uri:"id" binding:"min=1,required"`
}

// getOrderByID обрабатывает запрос на получение детальной информации о заказе
func (server *Server) getOrderByID(ctx *gin.Context) {
	var req getOrderByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем данные пользователя из токена
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	// Получаем заказ из базы данных
	order, err := server.store.GetOrderByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("заказ не найден")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Проверяем, имеет ли пользователь доступ к этому заказу
	// Клиент может видеть только свои заказы, провайдер - только заказы в своих категориях
	if user.ID != order.ClientID && user.Role.Role != sqlc.RoleAdmin {
		// Если пользователь не клиент и не администратор, проверяем, является ли он провайдером
		if user.Role.Role == sqlc.RoleProvider {
			// Проверяем, есть ли у провайдера услуги в категории заказа
			hasServices, err := server.checkProviderHasServicesInCategory(ctx, user.ID, order.CategoryID)
			if err != nil || !hasServices {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("у вас нет доступа к этому заказу")))
				return
			}
		} else {
			ctx.JSON(http.StatusForbidden, errorResponse(errors.New("у вас нет доступа к этому заказу")))
			return
		}
	}

	ctx.JSON(http.StatusOK, order)
}

// listOrdersRequest представляет запрос на получение списка заказов
type listOrdersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listOrdersRespons struct {
	Orders     []sqlc.ListOrdersByClientIDRow `json:"orders"`
	OrdersSize int                            `json:"orders_size"`
}

// listOrders обрабатывает запрос на получение списка заказов клиента
func (server *Server) listOrders(ctx *gin.Context) {
	var req listOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем данные пользователя из токена
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	// Создаем параметры для запроса
	arg := sqlc.ListOrdersByClientIDParams{
		ClientID: user.ID,
		Limit:    int64(req.PageSize),
		Offset:   int64((req.PageID - 1) * req.PageSize),
	}

	// Получаем заказы из базы данных
	orders, err := server.store.ListOrdersByClientID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Получаем кол-во записей
	orderSize, err := server.store.ListCountOrdersByClientID(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listOrdersRespons{
		Orders:     orders,
		OrdersSize: int(orderSize),
	}

	ctx.JSON(http.StatusOK, rsp)
}

// listAvailableOrdersRequest представляет запрос на получение доступных заказов для провайдера
type listAvailableOrdersRequest struct {
	CategoryFilter int64 `form:"category_filter"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAvailableOrders обрабатывает запрос на получение доступных заказов для провайдера
func (server *Server) listAvailableOrders(ctx *gin.Context) {
	var req listAvailableOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем данные пользователя из токена
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	// Проверяем, что пользователь - провайдер
	if user.Role.Role == sqlc.RoleClient {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("только провайдеры могут просматривать доступные заказы")))
		return
	}

	// Создаем параметры для запроса
	arg := sqlc.ListAvailableOrdersForProviderParams{
		// ProviderID: user.ID,
		Limit:      int64(req.PageSize),
		Offset:     int64((req.PageID - 1) * req.PageSize),
	}

	// Получаем заказы из базы данных
	orders, err := server.store.ListAvailableOrdersForProvider(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resultOrders := []sqlc.ListAvailableOrdersForProviderRow{}

	for _, order := range orders {
		if req.CategoryFilter != 0 {
			if order.ClientCity.String == user.City.String && order.CategoryID == req.CategoryFilter {
				resultOrders = append(resultOrders, order)
			}
		} else {
			if order.ClientCity.String == user.City.String {
				resultOrders = append(resultOrders, order)
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orders":      resultOrders,
		"orders_size": len(resultOrders),
	})
}

// acceptOrderRequest представляет запрос на принятие заказа провайдером
type acceptOrderRequest struct {
	ServiceID       int64  `json:"service_id" binding:"required"`       // ID услуги провайдера
	ProviderMessage string `json:"provider_message" binding:"required"` // Сообщение от провайдера
}

// acceptOrderResponse представляет ответ на принятие заказа
type acceptOrderResponse struct {
	OrderID         int64     `json:"order_id"`
	ClientID        int64     `json:"client_id"`
	ServiceID       int64     `json:"service_id"`
	ProviderID      int64     `json:"provider_id"`
	Status          string    `json:"status"`
	ProviderMessage string    `json:"provider_message"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// acceptOrder обрабатывает запрос на принятие заказа провайдером
func (server *Server) acceptOrder(ctx *gin.Context) {
	// Получаем ID заказа из параметров пути
	orderIDStr := ctx.Param("id")
	if orderIDStr == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("не указан ID заказа")))
		return
	}

	// Преобразуем ID заказа в int64
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("неверный формат ID заказа")))
		return
	}

	// Получаем данные из JSON тела запроса
	var req acceptOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем данные пользователя из токена
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	// Проверяем, что пользователь - провайдер
	if user.Role.Role != sqlc.RoleProvider && user.Role.Role != sqlc.RoleAdmin {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("только провайдеры могут принимать заказы")))
		return
	}

	// Проверяем, принадлежит ли услуга провайдеру
	service, err := server.store.GetServiceByID(ctx, req.ServiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("услуга не найдена")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if service.ProviderID != user.ID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("вы можете принимать заказы только на свои услуги")))
		return
	}

	// Создаем параметры для запроса
	arg := sqlc.AcceptOrderByProviderIDParams{
		Column1: sql.NullInt64{Int64: orderID, Valid: true},
		Column2: sql.NullInt64{Int64: user.ID, Valid: true},
		Column3: sql.NullInt64{Int64: req.ServiceID, Valid: true},
		Column4: sql.NullString{String: req.ProviderMessage, Valid: true},
	}

	// Обновляем заказ в базе данных
	updatedOrder, err := server.store.AcceptOrderByProviderID(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			// Проверяем причину, почему заказ не может быть принят
			order, checkErr := server.store.GetOrderByID(ctx, orderID)
			if checkErr != nil {
				ctx.JSON(http.StatusNotFound, errorResponse(errors.New("заказ не найден")))
				return
			}

			// Проверяем, не принят ли уже заказ
			if order.ProviderAccepted.Bool {
				ctx.JSON(http.StatusConflict, errorResponse(errors.New("этот заказ уже принят другим провайдером")))
				return
			}

			// Проверяем статус заказа
			if order.Status.StatusOrders != sqlc.StatusOrdersPending {
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("можно принять только заказ со статусом 'в ожидании'")))
				return
			}

			// Проверяем категорию заказа
			hasServices, err := server.checkProviderHasServicesInCategory(ctx, user.ID, order.CategoryID)
			if err != nil || !hasServices {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("вы не можете принять этот заказ, так как не предлагаете услуги в этой категории")))
				return
			}

			// Проверяем соответствие городов
			client, err := server.store.GetUserByID(ctx, order.ClientID)
			if err == nil && client.City.Valid {
				serviceObj, err := server.store.GetServiceByID(ctx, req.ServiceID)
				if err == nil && serviceObj.City.Valid {
					if serviceObj.City.String != client.City.String {
						ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("ваша услуга недоступна в городе клиента")))
						return
					}
				}
			}

			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("не удалось принять заказ")))
			return
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("указана несуществующая услуга")))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Формируем ответ
	rsp := acceptOrderResponse{
		OrderID:         updatedOrder.ID,
		ClientID:        updatedOrder.ClientID,
		ServiceID:       updatedOrder.ServiceID.Int64,
		ProviderID:      user.ID,
		Status:          string(updatedOrder.Status.StatusOrders),
		ProviderMessage: updatedOrder.ProviderMessage.String,
		UpdatedAt:       updatedOrder.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// updateOrderStatusRequest представляет запрос на обновление статуса заказа
type updateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending accepted completed cancelled"` // Новый статус заказа
}

// updateOrderStatus обрабатывает запрос на обновление статуса заказа
func (server *Server) updateOrderStatus(ctx *gin.Context) {
	// Получаем ID заказа из параметров пути
	orderIDStr := ctx.Param("id")
	if orderIDStr == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("не указан ID заказа")))
		return
	}

	// Преобразуем ID заказа в int64
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("неверный формат ID заказа")))
		return
	}

	// Получаем данные из JSON тела запроса
	var req updateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем данные пользователя из токена
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	// Получаем данные о заказе
	order, err := server.store.GetOrderByID(ctx, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("заказ не найден")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Проверяем права на обновление статуса
	// Клиент может отменить заказ, если он еще не принят, или отметить как выполненный
	// Провайдер может отметить заказ как выполненный или отмененный, только если он его принял
	if user.Role.Role != sqlc.RoleAdmin {
		if user.ID == order.ClientID {
			// Клиент может только отменить заказ или отметить как выполненный
			if req.Status == string(sqlc.StatusOrdersAccepted) {
				if string(order.Status.StatusOrders) == string(sqlc.StatusOrdersAccepted) {
					ctx.JSON(http.StatusForbidden, errorResponse(errors.New("нельзя отменить принятый заказ")))
					return
				}
			} else if req.Status == string(sqlc.StatusOrdersCancelled) {
				if string(order.Status.StatusOrders) != string(sqlc.StatusOrdersAccepted) {
					ctx.JSON(http.StatusForbidden, errorResponse(errors.New("можно отметить как выполненный только принятый заказ")))
					return
				}
			} else {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("клиент может только отменить заказ или отметить как выполненный")))
				return
			}
		} else if user.Role.Role == sqlc.RoleProvider {
			// Провайдер может изменить статус только своих принятых заказов
			if !order.ServiceID.Valid {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("заказ еще не принят провайдером")))
				return
			}

			// Проверяем, принадлежит ли заказ этому провайдеру
			service, err := server.store.GetServiceByID(ctx, order.ServiceID.Int64)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			if service.ProviderID != user.ID {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("вы можете изменять статус только своих заказов")))
				return
			}

			// Провайдер может только отметить заказ как выполненный или отмененный
			if req.Status != "completed" && req.Status != "cancelled" {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("провайдер может только отметить заказ как выполненный или отмененный")))
				return
			}
		} else {
			ctx.JSON(http.StatusForbidden, errorResponse(errors.New("у вас нет прав на обновление статуса этого заказа")))
			return
		}
	}

	// Создаем параметры для запроса
	var statusEnum sqlc.StatusOrders
	switch req.Status {
	case "pending":
		statusEnum = sqlc.StatusOrdersPending
	case "accepted":
		statusEnum = sqlc.StatusOrdersAccepted
	case "completed":
		statusEnum = sqlc.StatusOrdersCompleted
	case "cancelled":
		statusEnum = sqlc.StatusOrdersCancelled
	}

	arg := sqlc.UpdateOrderStatusParams{
		ID:     orderID,
		Status: sqlc.NullStatusOrders{StatusOrders: statusEnum, Valid: true},
	}

	// Обновляем статус заказа в базе данных
	updatedOrder, err := server.store.UpdateOrderStatus(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedOrder)
}

type updateOrderRequest struct {
	OrderID       int64  `json:"order_id" binding:"min=1,required"`
	CategoryID    int64  `json:"category_id" binding:"min=1,required"`
	ClientMessage string `json:"client_message" binding:"required"`
}

func (server *Server) updatedOrder(ctx *gin.Context) {
	var req updateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	isOrder, err := server.store.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.UpdateOrderParams{
		ID:            isOrder.ID,
		CategoryID:    req.CategoryID,
		ClientMessage: sql.NullString{String: req.ClientMessage, Valid: true},
	}

	updateOrder, err := server.store.UpdateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updateOrder)
}

// getOrdersStatisticsResponse представляет ответ со статистикой заказов
type getOrdersStatisticsResponse struct {
	PendingCount   int64 `json:"pending_count"`
	AcceptedCount  int64 `json:"accepted_count"`
	CompletedCount int64 `json:"completed_count"`
	CancelledCount int64 `json:"cancelled_count"`
	TotalCount     int64 `json:"total_count"`
}

// getOrdersStatistics обрабатывает запрос на получение статистики заказов для провайдера
func (server *Server) getOrdersStatistics(ctx *gin.Context) {
	// Получаем данные пользователя из токена
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	// Проверяем, что пользователь - провайдер
	if user.Role.Role != sqlc.RoleProvider && user.Role.Role != sqlc.RoleAdmin {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("только провайдеры могут просматривать статистику заказов")))
		return
	}

	// Получаем статистику заказов из базы данных
	stats, err := server.store.GetOrderStatistics(ctx, sql.NullInt64{Int64: user.ID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Формируем ответ
	rsp := getOrdersStatisticsResponse{
		PendingCount:   stats.PendingCount,
		AcceptedCount:  stats.AcceptedCount,
		CompletedCount: stats.CompletedCount,
		CancelledCount: stats.CancelledCount,
		TotalCount:     stats.TotalCount,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// getOrdersByCategoryRequest представляет запрос на получение заказов по категории
type getOrdersByCategoryRequest struct {
	CategoryID int64 `uri:"category_id" binding:"required"`
	PageID     int32 `form:"page_id" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// getOrdersByCategory обрабатывает запрос на получение заказов по категории
func (server *Server) getOrdersByCategory(ctx *gin.Context) {
	var req getOrdersByCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем данные пользователя из токена
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	// Проверяем права доступа (провайдер может просматривать заказы только в своих категориях)
	if user.Role.Role == sqlc.RoleProvider {
		hasServices, err := server.checkProviderHasServicesInCategory(ctx, user.ID, req.CategoryID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		if !hasServices {
			ctx.JSON(http.StatusForbidden, errorResponse(errors.New("вы можете просматривать заказы только в категориях ваших услуг")))
			return
		}
	} else if user.Role.Role != sqlc.RoleAdmin {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("только провайдеры и администраторы могут просматривать заказы по категориям")))
		return
	}

	// Создаем параметры для запроса
	arg := sqlc.GetOrdersByCategoryParams{
		CategoryID: req.CategoryID,
		Limit:      int64(req.PageSize),
		Offset:     int64((req.PageID - 1) * req.PageSize),
	}

	// Получаем заказы из базы данных
	orders, err := server.store.GetOrdersByCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (server *Server) deleteOrder(ctx *gin.Context) {
	var req getOrderByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteOrder(ctx, req.ID)
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

// Вспомогательная функция для проверки, имеет ли провайдер услуги в определенной категории
func (server *Server) checkProviderHasServicesInCategory(ctx *gin.Context, providerID, categoryID int64) (bool, error) {
	services, err := server.store.ListServicesByProviderIDAndCategory(ctx, sqlc.ListServicesByProviderIDAndCategoryParams{
		ProviderID: providerID,
		CategoryID: categoryID,
	})
	if err != nil {
		return false, err
	}
	return len(services) > 0, nil
}
