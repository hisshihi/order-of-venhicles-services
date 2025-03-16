package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

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

type User struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Country  *string `json:"country"`
	City     *string `json:"city"`
	District *string `json:"district"`
	Phone    string  `json:"phone"`
	Whatsapp string  `json:"whatsapp"`
}

type createUserResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

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

	accessToken, err := server.maker.CreateToken(user.ID, user.Email, string(user.Role.Role), server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		User: User{
			Username: user.Username,
			Email:    user.Email,
			Country:  &user.Country.String,
			City:     &user.City.String,
			District: &user.District.String,
			Phone:    user.Phone,
			Whatsapp: user.Whatsapp,
		},
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password_hash" binding:"min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Если запрос верный, находим пользователя
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Проверяем, правильно ли пользователь ввёл пароль
	err = util.CheckPassword(req.PasswordHash, user.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.maker.CreateToken(user.ID, user.Email, string(user.Role.Role), server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getCurrentUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Country  string `json:"country,omitempty"`
	City     string `json:"city,omitempty"`
	District string `json:"district,omitempty"`
	Phone    string `json:"phone"`
	Whatsapp string `json:"whatsapp"`
}

func (server *Server) getCurrentUser(ctx *gin.Context) {
	// Получаем payload из токена авторизации
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		err := errors.New("требуется авторизация")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Приводим payload к нужному типу
	tokenPayload, ok := payload.(*util.Payload)
	if !ok {
		err := errors.New("неверный тип payload")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Получаем пользователя по ID из токена
	user, err := server.store.GetUserByIDFromUser(ctx, tokenPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Создаем безопасный ответ без хеша пароля
	rsp := getCurrentUserResponse{
		Email:    user.Email,
		Username: user.Username,
		Country:  user.Country.String,
		City:     user.City.String,
		District: user.District.String,
		Phone:    user.Phone,
		Whatsapp: user.Whatsapp,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// getUserByID - метод для администраторов, позволяющий получить данные любого пользователя по ID
func (server *Server) getUserByID(ctx *gin.Context) {
	var req struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем пользователя по ID из URI
	user, err := server.store.GetUserByIDFromUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Создаем безопасный ответ без хеша пароля
	userResponse := struct {
		ID        int64     `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Country   string    `json:"country,omitempty"`
		City      string    `json:"city,omitempty"`
		District  string    `json:"district,omitempty"`
		Phone     string    `json:"phone"`
		Whatsapp  string    `json:"whatsapp"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Country:   user.Country.String,
		City:      user.City.String,
		District:  user.District.String,
		Phone:     user.Phone,
		Whatsapp:  user.Whatsapp,
		Role:      string(user.Role.Role),
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, userResponse)
}
