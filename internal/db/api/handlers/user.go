package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/hisshihi/order-of-venhicles-services/pkg/util"
	"github.com/lib/pq"
)

type createUserParams struct {
	Username     string  `json:"username" binding:"required,alphanum"`
	Email        string  `json:"email" binding:"required,email"`
	PasswordHash string  `json:"password_hash" binding:"min=6"`
	Country      *string `json:"country,omitempty"`
	City         *string `json:"city,omitempty"`
	District     *string `json:"district,omitempty"`
	Phone        string  `json:"phone" binding:"required"`
	Whatsapp     string  `json:"whatsapp" binding:"required"`
	Role         *string `json:"role,omitempty"`
}

type createUserResponse struct {
	Username     string  `json:"username" binding:"required,alphanum"`
	Email        string  `json:"email" binding:"required,email"`
	Country      *string `json:"country,omitempty"`
	City         *string `json:"city,omitempty"`
	District     *string `json:"district,omitempty"`
	Phone        string  `json:"phone" binding:"required"`
	Whatsapp     string  `json:"whatsapp" binding:"required"`
}

// TODO: При создании добавить возвращение токена

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Безопасность: проверка, чтобы пользователь не мог сам себе назначить роль admin
	if req.Role != nil && *req.Role == string(sqlc.RoleAdmin) {
		err := errors.New("установка роли администратора при регистрации запрещена")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Подготовка параметров с обработкой необязательных полей
	arg := sqlc.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Whatsapp:     req.Whatsapp,
		// По умолчанию все пользователи создаются как клиенты
		Role: sqlc.NullRole{Role: sqlc.RoleClient, Valid: true},
	}

	// Обработка необязательных полей
	if req.Country != nil {
		arg.Country = sql.NullString{String: *req.Country, Valid: true}
	}

	if req.City != nil {
		arg.City = sql.NullString{String: *req.City, Valid: true}
	}

	if req.District != nil {
		arg.District = sql.NullString{String: *req.District, Valid: true}
	}

	// Если передана роль, используем её вместо значения по умолчанию
	// Роли ограничены только client, provider или partner
	if req.Role != nil {
		roleValue := *req.Role
		// Проверка допустимости роли
		switch roleValue {
		case string(sqlc.RoleClient), string(sqlc.RoleProvider), string(sqlc.RolePartner):
			arg.Role = sqlc.NullRole{Role: sqlc.Role(roleValue), Valid: true}
		default:
			// Если роль не поддерживается, используем роль по умолчанию
			// Здесь вместо тихого игнорирования лучше вернуть ошибку
			err := errors.New("недопустимая роль: разрешены только client, provider, partner")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	user, err := server.store.CreateUser(ctx, arg)
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

	rsp := createUserResponse{
		Username: user.Username,
		Email: user.Email,
		Country: &user.Country.String,
		City: &user.City.String,
		District: &user.District.String,
		Phone: user.Phone,
		Whatsapp: user.Whatsapp,
	}

	ctx.JSON(http.StatusOK, rsp)
}
