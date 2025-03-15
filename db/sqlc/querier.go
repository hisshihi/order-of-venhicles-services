// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreatePromoCode(ctx context.Context, arg CreatePromoCodeParams) (PromoCode, error)
	CreateService(ctx context.Context, arg CreateServiceParams) (Service, error)
	CreateServiceCategory(ctx context.Context, arg CreateServiceCategoryParams) (ServiceCategory, error)
	CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteOrder(ctx context.Context, id int64) error
	DeletePromoCode(ctx context.Context, id int64) error
	DeleteService(ctx context.Context, id int64) error
	DeleteServiceCategory(ctx context.Context, id int64) error
	DeleteSubscription(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetOrderByID(ctx context.Context, id int64) (Order, error)
	GetPromoCodeByID(ctx context.Context, id int64) (PromoCode, error)
	GetPromoCodeByPartnerID(ctx context.Context, partnerID int64) (PromoCode, error)
	GetServiceByID(ctx context.Context, id int64) (Service, error)
	GetServiceCategoryByID(ctx context.Context, id int64) (ServiceCategory, error)
	GetSubscriptionByID(ctx context.Context, id int64) (Subscription, error)
	GetUserByIDFromAdmin(ctx context.Context, id int64) (User, error)
	GetUserByIDFromUser(ctx context.Context, id int64) (interface{}, error)
	ListOrders(ctx context.Context, arg ListOrdersParams) ([]Order, error)
	ListPromoCodes(ctx context.Context, arg ListPromoCodesParams) ([]PromoCode, error)
	ListServiceCategories(ctx context.Context, arg ListServiceCategoriesParams) ([]ServiceCategory, error)
	ListServices(ctx context.Context, arg ListServicesParams) ([]Service, error)
	ListServicesByCategory(ctx context.Context, categoryID int64) ([]Service, error)
	ListServicesByTitle(ctx context.Context, dollar_1 sql.NullString) ([]Service, error)
	ListSubscriptions(ctx context.Context, arg ListSubscriptionsParams) ([]Subscription, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	ListUsersByEmail(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	ListUsersByRole(ctx context.Context, role NullRole) ([]User, error)
	ListUsersByUsername(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdatePromoCode(ctx context.Context, arg UpdatePromoCodeParams) (PromoCode, error)
	UpdateService(ctx context.Context, arg UpdateServiceParams) (Service, error)
	UpdateServiceCategory(ctx context.Context, arg UpdateServiceCategoryParams) (ServiceCategory, error)
	UpdateSubscription(ctx context.Context, arg UpdateSubscriptionParams) (Subscription, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
