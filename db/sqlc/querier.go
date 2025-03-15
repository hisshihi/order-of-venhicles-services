// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateServiceCategory(ctx context.Context, arg CreateServiceCategoryParams) (ServiceCategory, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteServiceCategory(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetServiceCategoryByID(ctx context.Context, id int64) (ServiceCategory, error)
	GetUserByIDFromAdmin(ctx context.Context, id int64) (User, error)
	GetUserByIDFromUser(ctx context.Context, id int64) (interface{}, error)
	ListServiceCategories(ctx context.Context, arg ListServiceCategoriesParams) ([]ServiceCategory, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	ListUsersByEmail(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	ListUsersByRole(ctx context.Context, role NullRole) ([]User, error)
	ListUsersByUsername(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	UpdateServiceCategory(ctx context.Context, arg UpdateServiceCategoryParams) (ServiceCategory, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
