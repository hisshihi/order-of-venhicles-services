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
	CheckAndUpdateExpiredSubscriptions(ctx context.Context) ([]Subscription, error)
	// Проверяет, оставил ли клиент отзыв по данному заказу
	CheckIfClientReviewedOrder(ctx context.Context, arg CheckIfClientReviewedOrderParams) (bool, error)
	// Проверяет, добавлен ли услугодатель в избранное клиента
	CheckIfProviderIsFavorite(ctx context.Context, arg CheckIfProviderIsFavoriteParams) (bool, error)
	// Создает новое сообщение
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error)
	CreatePromoCode(ctx context.Context, arg CreatePromoCodeParams) (PromoCode, error)
	// Создает новый отзыв от клиента об услугодателе
	CreateReview(ctx context.Context, arg CreateReviewParams) (Review, error)
	CreateService(ctx context.Context, arg CreateServiceParams) (Service, error)
	CreateServiceCategory(ctx context.Context, name string) (ServiceCategory, error)
	CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteOrder(ctx context.Context, id int64) error
	DeletePayment(ctx context.Context, id int64) error
	DeletePromoCode(ctx context.Context, id int64) error
	// Удаляет отзыв (только если пользователь является автором или администратором)
	DeleteReview(ctx context.Context, arg DeleteReviewParams) error
	DeleteService(ctx context.Context, arg DeleteServiceParams) error
	DeleteServiceCategory(ctx context.Context, id int64) error
	DeleteSubscription(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetActiveSubscriptionForProvider(ctx context.Context, providerID int64) (Subscription, error)
	// Получает всех поставщиков, которые использовали данный промокод
	GetAllProvidersByPartnerPromos(ctx context.Context, arg GetAllProvidersByPartnerPromosParams) ([]GetAllProvidersByPartnerPromosRow, error)
	// Получает среднюю оценку услугодателя
	GetAverageRatingForProvider(ctx context.Context, providerID int64) (GetAverageRatingForProviderRow, error)
	// Получает историю переписки между двумя пользователями
	GetMessagesByUsers(ctx context.Context, arg GetMessagesByUsersParams) ([]GetMessagesByUsersRow, error)
	GetOrderByID(ctx context.Context, id int64) (GetOrderByIDRow, error)
	// Получает статистику заказов для услугодателя
	GetOrderStatistics(ctx context.Context, dollar_1 sql.NullInt64) (GetOrderStatisticsRow, error)
	// Получает список заказов по категории
	GetOrdersByCategory(ctx context.Context, arg GetOrdersByCategoryParams) ([]GetOrdersByCategoryRow, error)
	GetPaymentByID(ctx context.Context, id int64) (Payment, error)
	GetPromoCodeByID(ctx context.Context, id int64) (PromoCode, error)
	GetPromoCodeByPartnerID(ctx context.Context, partnerID int64) (PromoCode, error)
	GetProvidersByPromoCode(ctx context.Context, arg GetProvidersByPromoCodeParams) ([]GetProvidersByPromoCodeRow, error)
	// Получает конкретный отзыв по ID
	GetReviewByID(ctx context.Context, id int64) (GetReviewByIDRow, error)
	// Получает все отзывы об услугодателе
	GetReviewsByProviderID(ctx context.Context, arg GetReviewsByProviderIDParams) ([]GetReviewsByProviderIDRow, error)
	GetServiceByID(ctx context.Context, id int64) (GetServiceByIDRow, error)
	GetServiceCategoryByID(ctx context.Context, id int64) (ServiceCategory, error)
	GetServicesByProviderID(ctx context.Context, arg GetServicesByProviderIDParams) ([]GetServicesByProviderIDRow, error)
	GetSubscriptionByID(ctx context.Context, id int64) (Subscription, error)
	GetSubscriptionByProviderID(ctx context.Context, providerID int64) (Subscription, error)
	// Получает количество непрочитанных сообщений для пользователя
	GetUnreadMessagesCount(ctx context.Context, receiverID int64) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByIDFromAdmin(ctx context.Context, id int64) (User, error)
	GetUserByIDFromUser(ctx context.Context, id int64) (User, error)
	// Получает список недавних чатов пользователя
	GetUserRecentChats(ctx context.Context, arg GetUserRecentChatsParams) ([]GetUserRecentChatsRow, error)
	// Получает список доступных заказов для провайдера услуг
	ListAvailableOrdersForProvider(ctx context.Context, arg ListAvailableOrdersForProviderParams) ([]ListAvailableOrdersForProviderRow, error)
	// Получает список избранных услугодателей клиента
	ListFavoriteProviders(ctx context.Context, arg ListFavoriteProvidersParams) ([]ListFavoriteProvidersRow, error)
	ListOrdersByClientID(ctx context.Context, arg ListOrdersByClientIDParams) ([]ListOrdersByClientIDRow, error)
	ListPayments(ctx context.Context, arg ListPaymentsParams) ([]Payment, error)
	ListPromoCodes(ctx context.Context, arg ListPromoCodesParams) ([]PromoCode, error)
	ListPromoCodesByPartnerID(ctx context.Context, arg ListPromoCodesByPartnerIDParams) ([]PromoCode, error)
	ListServiceCategories(ctx context.Context, arg ListServiceCategoriesParams) ([]ServiceCategory, error)
	ListServices(ctx context.Context, arg ListServicesParams) ([]ListServicesRow, error)
	ListServicesByCategory(ctx context.Context, arg ListServicesByCategoryParams) ([]ListServicesByCategoryRow, error)
	ListServicesByLocation(ctx context.Context, arg ListServicesByLocationParams) ([]ListServicesByLocationRow, error)
	ListServicesByProviderIDAndCategory(ctx context.Context, arg ListServicesByProviderIDAndCategoryParams) ([]ListServicesByProviderIDAndCategoryRow, error)
	ListSubscriptions(ctx context.Context, arg ListSubscriptionsParams) ([]Subscription, error)
	ListSubscriptionsByProviderID(ctx context.Context, arg ListSubscriptionsByProviderIDParams) ([]Subscription, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	ListUsersByEmail(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	ListUsersByRole(ctx context.Context, role NullRole) ([]User, error)
	ListUsersByUsername(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	// Отмечает сообщения как прочитанные
	MarkMessagesAsRead(ctx context.Context, arg MarkMessagesAsReadParams) error
	// Удаляет услугодателя из избранного клиента
	RemoveProviderFromFavorites(ctx context.Context, arg RemoveProviderFromFavoritesParams) error
	SearchServices(ctx context.Context, arg SearchServicesParams) ([]SearchServicesRow, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) (Order, error)
	UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (Payment, error)
	UpdatePromoCode(ctx context.Context, arg UpdatePromoCodeParams) (PromoCode, error)
	UpdateService(ctx context.Context, arg UpdateServiceParams) (Service, error)
	UpdateServiceCategory(ctx context.Context, arg UpdateServiceCategoryParams) (ServiceCategory, error)
	UpdateSubscription(ctx context.Context, arg UpdateSubscriptionParams) (Subscription, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
