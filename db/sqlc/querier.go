// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"context"
	"database/sql"
)

type Querier interface {
	// Провайдер принимает заказ и указывает свою услугу
	AcceptOrderByProviderID(ctx context.Context, arg AcceptOrderByProviderIDParams) (Order, error)
	// Добавляет услугодателя в избранное клиента
	AddProviderToFavorites(ctx context.Context, arg AddProviderToFavoritesParams) (Favorite, error)
	BlockedUser(ctx context.Context, id int64) (sql.NullBool, error)
	ChangePassword(ctx context.Context, arg ChangePasswordParams) error
	CheckAndUpdateExpiredSubscriptions(ctx context.Context) ([]Subscription, error)
	// Проверяет, оставил ли клиент отзыв по данному заказу
	CheckIfClientReviewedOrder(ctx context.Context, arg CheckIfClientReviewedOrderParams) (bool, error)
	// Проверяет, добавлен ли услугодатель в избранное клиента
	CheckIfProviderIsFavorite(ctx context.Context, arg CheckIfProviderIsFavoriteParams) (bool, error)
	// Подсчитывает количество откликов на конкретный заказ
	// Входные параметры: ID заказа
	// Возвращает: количество откликов (число)
	CountOrderResponsesByOrderId(ctx context.Context, orderID int64) (int64, error)
	CountOrders(ctx context.Context) (int64, error)
	CountPromoCode(ctx context.Context) (int64, error)
	CountReviews(ctx context.Context) (int64, error)
	CountService(ctx context.Context) (int64, error)
	CountServicesByProviderID(ctx context.Context, providerID int64) (int64, error)
	CountSubscriptions(ctx context.Context) (int64, error)
	CountSupportMessages(ctx context.Context) (int64, error)
	CountUsers(ctx context.Context) (int64, error)
	CreateCity(ctx context.Context, name string) (City, error)
	// Создает новое сообщение
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	// Создает новый отклик на заказ от провайдера услуг
	// Входные параметры: ID заказа, ID провайдера, сообщение, предложенная цена, выбран ли отклик
	// Возвращает: созданную запись отклика
	CreateOrderResponse(ctx context.Context, arg CreateOrderResponseParams) (OrderResponse, error)
	CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error)
	CreatePendingSubscription(ctx context.Context, arg CreatePendingSubscriptionParams) (PendingSubscription, error)
	CreatePromoCode(ctx context.Context, arg CreatePromoCodeParams) (PromoCode, error)
	// Создает новый отзыв от клиента об услугодателе
	CreateReview(ctx context.Context, arg CreateReviewParams) (Review, error)
	CreateService(ctx context.Context, arg CreateServiceParams) (Service, error)
	CreateServiceCategory(ctx context.Context, arg CreateServiceCategoryParams) (ServiceCategory, error)
	CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error)
	CreateSubtitle(ctx context.Context, name string) (SubtitleCategory, error)
	CreateSupportMessage(ctx context.Context, arg CreateSupportMessageParams) (SupportMessage, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteCity(ctx context.Context, id int64) error
	DeleteOrder(ctx context.Context, id int64) error
	// Удаляет отклик по его ID
	// Входные параметры: ID отклика
	// Возвращает: ничего
	DeleteOrderResponse(ctx context.Context, id int64) error
	DeleteOrdersByCategoryID(ctx context.Context, categoryID int64) (int64, error)
	DeleteOrdersBySubcategoryID(ctx context.Context, subtitleCategoryID sql.NullInt64) (int64, error)
	DeletePayment(ctx context.Context, id int64) error
	DeletePendingSubscriptionByPaymentID(ctx context.Context, paymentID int64) error
	DeletePromoCode(ctx context.Context, id int64) error
	// Удаляет отзыв (только если пользователь является автором или администратором)
	DeleteReview(ctx context.Context, arg DeleteReviewParams) error
	DeleteService(ctx context.Context, arg DeleteServiceParams) error
	DeleteServiceCategory(ctx context.Context, id int64) (int64, error)
	DeleteServicesByCategoryID(ctx context.Context, categoryID int64) (int64, error)
	DeleteServicesBySubcategoryID(ctx context.Context, subtitleCategoryID sql.NullInt64) (int64, error)
	DeleteSubscription(ctx context.Context, id int64) error
	DeleteSubtitleCategory(ctx context.Context, id int64) (int64, error)
	DeleteSupportMessage(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	// Фильтрация услуг по цене
	FilterServiceByPrice(ctx context.Context, arg FilterServiceByPriceParams) ([]Service, error)
	GetActiveSubscriptionForProvider(ctx context.Context, providerID int64) (Subscription, error)
	// Получает всех поставщиков, которые использовали данный промокод
	GetAllProvidersByPartnerPromos(ctx context.Context, arg GetAllProvidersByPartnerPromosParams) ([]GetAllProvidersByPartnerPromosRow, error)
	// Получает среднюю оценку услугодателя
	GetAverageRatingForProvider(ctx context.Context, providerID int64) (GetAverageRatingForProviderRow, error)
	GetBlockerUser(ctx context.Context, id int64) (GetBlockerUserRow, error)
	GetCityByID(ctx context.Context, id int64) (City, error)
	// Получает историю переписки между двумя пользователями
	GetMessagesByUsers(ctx context.Context, arg GetMessagesByUsersParams) ([]GetMessagesByUsersRow, error)
	GetOrderByID(ctx context.Context, id int64) (GetOrderByIDRow, error)
	// Получает отклик по его ID
	// Входные параметры: ID отклика
	// Возвращает: запись отклика или ничего, если отклик не найден
	GetOrderResponseById(ctx context.Context, id int64) (OrderResponse, error)
	// Проверяет наличие отклика от конкретного провайдера на конкретный заказ
	// Входные параметры: ID заказа, ID провайдера
	// Возвращает: запись отклика или ничего, если отклик не найден
	GetOrderResponseByOrderAndProvider(ctx context.Context, arg GetOrderResponseByOrderAndProviderParams) (OrderResponse, error)
	// Получает все отклики для конкретного заказа с информацией о провайдерах
	// Входные параметры: ID заказа
	// Возвращает: список откликов с данными провайдеров (имя, телефон, whatsapp, фото)
	GetOrderResponsesByOrderId(ctx context.Context, orderID int64) ([]GetOrderResponsesByOrderIdRow, error)
	// Получает все отклики, созданные конкретным провайдером, с информацией о связанных заказах
	// Входные параметры: ID провайдера, лимит и смещение для пагинации
	// Возвращает: список откликов с основной информацией о заказах
	GetOrderResponsesByProviderId(ctx context.Context, arg GetOrderResponsesByProviderIdParams) ([]GetOrderResponsesByProviderIdRow, error)
	// Получает статистику заказов для услугодателя
	GetOrderStatistics(ctx context.Context, dollar_1 sql.NullInt64) (GetOrderStatisticsRow, error)
	// Получает список заказов по категории
	GetOrdersByCategory(ctx context.Context, arg GetOrdersByCategoryParams) ([]GetOrdersByCategoryRow, error)
	// Получает список заказов по подкатегориям
	GetOrdersBySubCategory(ctx context.Context, arg GetOrdersBySubCategoryParams) ([]GetOrdersBySubCategoryRow, error)
	// Получает все заказы, на которые провайдер оставил отклики (даже если не выбран)
	// Входные параметры: ID провайдера, лимит и смещение для пагинации
	// Возвращает: список заказов с базовой информацией о клиенте и категории
	GetOrdersWithResponsesByProvider(ctx context.Context, arg GetOrdersWithResponsesByProviderParams) ([]GetOrdersWithResponsesByProviderRow, error)
	// Получает список заказов, для которых выбран провайдер, с детальной информацией
	// Входные параметры: лимит и смещение для пагинации
	// Возвращает: список заказов с данными о клиенте, категории и провайдере
	GetOrdersWithSelectedProvider(ctx context.Context, arg GetOrdersWithSelectedProviderParams) ([]GetOrdersWithSelectedProviderRow, error)
	GetPaymentByID(ctx context.Context, id int64) (Payment, error)
	GetPendingSubscriptionByPaymentID(ctx context.Context, paymentID int64) (PendingSubscription, error)
	GetPromoCodeByCode(ctx context.Context, code string) (PromoCode, error)
	GetPromoCodeByID(ctx context.Context, id int64) (PromoCode, error)
	GetPromoCodeByPartnerID(ctx context.Context, partnerID int64) (PromoCode, error)
	GetProvidersByPromoCode(ctx context.Context, arg GetProvidersByPromoCodeParams) ([]GetProvidersByPromoCodeRow, error)
	// Получает конкретный отзыв по ID
	GetReviewByID(ctx context.Context, id int64) (GetReviewByIDRow, error)
	// Получает все отзывы об услугодателе
	GetReviewsByProviderID(ctx context.Context, arg GetReviewsByProviderIDParams) ([]GetReviewsByProviderIDRow, error)
	// Получает данные выбранного провайдера для конкретного заказа
	// Входные параметры: ID заказа
	// Возвращает: запись пользователя-провайдера или ничего, если провайдер не выбран
	GetSelectedProviderForOrder(ctx context.Context, id int64) (User, error)
	// Получает выбранный отклик для конкретного заказа
	// Входные параметры: ID заказа
	// Возвращает: запись выбранного отклика или ничего, если отклик не выбран
	GetSelectedResponseForOrder(ctx context.Context, orderID int64) (OrderResponse, error)
	GetServiceByID(ctx context.Context, id int64) (GetServiceByIDRow, error)
	GetServiceCategoryByID(ctx context.Context, id int64) (ServiceCategory, error)
	GetServiceCategoryBySlug(ctx context.Context, slug string) (ServiceCategory, error)
	GetServicesByProviderID(ctx context.Context, arg GetServicesByProviderIDParams) ([]GetServicesByProviderIDRow, error)
	GetSubscriptionByID(ctx context.Context, id int64) (Subscription, error)
	GetSubscriptionByProviderID(ctx context.Context, providerID int64) (Subscription, error)
	GetSubtitleCategoryByID(ctx context.Context, id int64) (SubtitleCategory, error)
	// Получает количество непрочитанных сообщений для пользователя
	GetUnreadMessagesCount(ctx context.Context, receiverID int64) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByIDFromAdmin(ctx context.Context, id int64) (User, error)
	GetUserByIDFromUser(ctx context.Context, id int64) (User, error)
	// Получает список недавних чатов пользователя
	GetUserRecentChats(ctx context.Context, arg GetUserRecentChatsParams) ([]GetUserRecentChatsRow, error)
	// Проверяет, откликнулся ли провайдер на заказ
	// Входные параметры: ID заказа, ID провайдера
	// Возвращает: булево значение (true/false)
	HasProviderRespondedToOrder(ctx context.Context, arg HasProviderRespondedToOrderParams) (bool, error)
	// Получает список доступных заказов для провайдера услуг
	ListAvailableOrdersForProvider(ctx context.Context, arg ListAvailableOrdersForProviderParams) ([]ListAvailableOrdersForProviderRow, error)
	ListBlockedUsers(ctx context.Context) ([]ListBlockedUsersRow, error)
	ListCity(ctx context.Context) ([]City, error)
	ListCountAvailableOrdersForProvider(ctx context.Context, providerID int64) (int64, error)
	ListCountOrdersByClientID(ctx context.Context, clientID int64) (int64, error)
	ListCountServicesByCatetegory(ctx context.Context, categoryID int64) (int64, error)
	// Получает список избранных услугодателей клиента
	ListFavoriteProviders(ctx context.Context, arg ListFavoriteProvidersParams) ([]ListFavoriteProvidersRow, error)
	ListOrders(ctx context.Context, arg ListOrdersParams) ([]Order, error)
	// Проверка на блокировку провайдера, если он есть
	ListOrdersByClientID(ctx context.Context, arg ListOrdersByClientIDParams) ([]ListOrdersByClientIDRow, error)
	ListPartners(ctx context.Context) ([]ListPartnersRow, error)
	ListPayments(ctx context.Context, arg ListPaymentsParams) ([]Payment, error)
	ListPromoCodes(ctx context.Context, arg ListPromoCodesParams) ([]PromoCode, error)
	ListPromoCodesByPartnerID(ctx context.Context, arg ListPromoCodesByPartnerIDParams) ([]PromoCode, error)
	ListReview(ctx context.Context, arg ListReviewParams) ([]Review, error)
	ListServiceCategories(ctx context.Context) ([]ServiceCategory, error)
	ListServices(ctx context.Context, arg ListServicesParams) ([]ListServicesRow, error)
	ListServicesByCategory(ctx context.Context, arg ListServicesByCategoryParams) ([]ListServicesByCategoryRow, error)
	ListServicesByLocation(ctx context.Context, arg ListServicesByLocationParams) ([]ListServicesByLocationRow, error)
	ListServicesByProviderIDAndCategory(ctx context.Context, arg ListServicesByProviderIDAndCategoryParams) ([]ListServicesByProviderIDAndCategoryRow, error)
	ListServicesByProviderIDAndSubCategory(ctx context.Context, arg ListServicesByProviderIDAndSubCategoryParams) ([]ListServicesByProviderIDAndSubCategoryRow, error)
	ListSubscriptions(ctx context.Context, arg ListSubscriptionsParams) ([]Subscription, error)
	ListSubscriptionsByProviderID(ctx context.Context, arg ListSubscriptionsByProviderIDParams) ([]Subscription, error)
	ListSubtitleCategory(ctx context.Context) ([]SubtitleCategory, error)
	ListSupportMessages(ctx context.Context, arg ListSupportMessagesParams) ([]SupportMessage, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	ListUsersByEmail(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	ListUsersByRole(ctx context.Context, role NullRole) ([]User, error)
	ListUsersByUsername(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	// Отмечает сообщения как прочитанные
	MarkMessagesAsRead(ctx context.Context, arg MarkMessagesAsReadParams) error
	// Удаляет услугодателя из избранного клиента
	RemoveProviderFromFavorites(ctx context.Context, arg RemoveProviderFromFavoritesParams) error
	SearchServices(ctx context.Context, arg SearchServicesParams) ([]SearchServicesRow, error)
	// Выбирает провайдера для выполнения заказа (транзакционный запрос)
	// Обновляет статус отклика на "выбран" и устанавливает выбранного провайдера в заказе
	// Входные параметры: ID отклика
	// Возвращает: обновленную запись заказа
	SelectProviderForOrder(ctx context.Context, dollar_1 sql.NullInt64) (Order, error)
	UnblockUser(ctx context.Context, id int64) (sql.NullBool, error)
	// Отменяет выбор провайдера для заказа (транзакционный запрос)
	// Сбрасывает статус отклика и очищает поле выбранного провайдера в заказе
	// Входные параметры: ID отклика
	// Возвращает: количество обновленных строк (должно быть 1)
	UnselectProviderForOrder(ctx context.Context, dollar_1 sql.NullInt64) (int64, error)
	UpdateCity(ctx context.Context, arg UpdateCityParams) (City, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	// Обновляет информацию в существующем отклике
	// Входные параметры: ID отклика, новое сообщение, новая предложенная цена
	// Возвращает: обновленную запись отклика
	UpdateOrderResponse(ctx context.Context, arg UpdateOrderResponseParams) (OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) (Order, error)
	UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (Payment, error)
	UpdatePaymentStatus(ctx context.Context, arg UpdatePaymentStatusParams) (Payment, error)
	UpdatePromoCodeByID(ctx context.Context, arg UpdatePromoCodeByIDParams) error
	UpdateService(ctx context.Context, arg UpdateServiceParams) (Service, error)
	UpdateServiceCategory(ctx context.Context, arg UpdateServiceCategoryParams) (ServiceCategory, error)
	UpdateSubscription(ctx context.Context, arg UpdateSubscriptionParams) (Subscription, error)
	UpdateSubtitleCategory(ctx context.Context, arg UpdateSubtitleCategoryParams) (SubtitleCategory, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserForAdmin(ctx context.Context, arg UpdateUserForAdminParams) (User, error)
}

var _ Querier = (*Queries)(nil)
