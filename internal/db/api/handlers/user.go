package api

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/hisshihi/order-of-venhicles-services/pkg/util"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username     string                `form:"username" binding:"required"`
	Email        string                `form:"email" binding:"required,email"`
	PasswordHash string                `form:"password_hash" binding:"min=6"`
	Country      *string               `form:"country,omitempty"`
	City         *string               `form:"city,omitempty"`
	District     *string               `form:"district,omitempty"`
	Phone        string                `form:"phone" binding:"required"`
	Whatsapp     string                `form:"whatsapp" binding:"required"`
	Role         *string               `form:"role,omitempty"`
	PhotoUrl     *multipart.FileHeader `form:"photo_url"`
}

type User struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Country   *string `json:"country"`
	City      *string `json:"city"`
	District  *string `json:"district"`
	Phone     string  `json:"phone"`
	Whatsapp  string  `json:"whatsapp"`
	PhotoUrl  string  `json:"photo_url,omitempty"`
	PhotoMime string  `json:"photo_mime,omitempty"`
}

type createUserResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

func (server *Server) createUser(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 1<<20+1024)

	var req createUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
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

	var photoBytes []byte
	if req.PhotoUrl != nil {
		if req.PhotoUrl.Size > 1<<29 {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("размер файла превышает 1МБ")))
			return
		}

		file, err := req.PhotoUrl.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		defer file.Close()

		// Считываем весь файл сразу
		photoBytes, err = io.ReadAll(file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Проверка типа контента (опционально)
		contentType := http.DetectContentType(photoBytes)
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
		}
		if !allowedTypes[contentType] {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("неподдерживаемый тип файла")))
			return
		}
	}

	// Подготовка параметров с обработкой необязательных полей
	arg := sqlc.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Whatsapp:     req.Whatsapp,
		// По умолчанию все пользователи создаются как клиенты
		Role:     sqlc.NullRole{Role: sqlc.RoleClient, Valid: true},
		Country:  sql.NullString{String: *req.Country, Valid: true},
		City:     sql.NullString{String: *req.City, Valid: true},
		District: sql.NullString{String: *req.District, Valid: true},
		PhotoUrl: photoBytes,
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

	contentType := http.DetectContentType(photoBytes)

	rsp := createUserResponse{
		User: User{
			Username:  user.Username,
			Email:     user.Email,
			Country:   &user.Country.String,
			City:      &user.City.String,
			District:  &user.District.String,
			Phone:     user.Phone,
			Whatsapp:  user.Whatsapp,
			PhotoUrl:  base64.StdEncoding.EncodeToString(user.PhotoUrl),
			PhotoMime: contentType,
		},
		AccessToken: accessToken,
	}

	ctx.SetCookie(
		"auth_token",
		accessToken,
		int(server.config.AccessTokenDuration.Seconds()),
		"/",
		"",
		false, // TODO: поменять на true для продакшена
		true,
	)

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password_hash" binding:"min=6"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        userResponse
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
		User: userResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     string(user.Role.Role),
		},
	}

	ctx.SetCookie(
		"auth_token",
		accessToken,
		int(server.config.AccessTokenDuration.Seconds()),
		"/",
		"",
		false, // TODO: поменять на true для продакшена
		true,
	)

	ctx.JSON(http.StatusOK, rsp)
}

// logoutUser обрабатывает запрос на выход из системы
func (server *Server) logoutUser(ctx *gin.Context) {
	// Удаляем cookie с токеном
	ctx.SetCookie(
		"auth_token",
		"", // пустой токен
		-1, // отрицательное время жизни для удаления
		"/",
		"",
		false,
		true,
	)

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Выход выполнен успешно"})
}

type getCurrentUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Country  string `json:"country,omitempty"`
	City     string `json:"city,omitempty"`
	District string `json:"district,omitempty"`
	Phone    string `json:"phone"`
	Whatsapp string `json:"whatsapp"`
	Role     string `json:"role"`
	PhotoUrl string `json:"photo_url,omitempty"`
}

func (server *Server) getCurrentUser(ctx *gin.Context) {
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
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
		Role:     string(user.Role.Role),
		PhotoUrl: base64.StdEncoding.EncodeToString(user.PhotoUrl),
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

// Профиль пользователя
func (server *Server) profileUser(ctx *gin.Context) {
	var req struct {
		ID       int64 `form:"id" binding:"required,min=1"`
		PageID   int32 `form:"page_id" binding:"required,min=1"`
		PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
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

	var averageRating sqlc.GetAverageRatingForProviderRow
	var allRating []sqlc.GetReviewsByProviderIDRow
	if user.Role.Role == sqlc.RoleProvider {
		averageRating, err = server.store.GetAverageRatingForProvider(ctx, user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		arg := sqlc.GetReviewsByProviderIDParams{
			ProviderID: user.ID,
			Limit:      int64(req.PageSize),
			Offset:     int64((req.PageID - 1) * req.PageSize),
		}

		allRating, err = server.store.GetReviewsByProviderID(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	userResponse := struct {
		ID            int64                               `json:"id"`
		Username      string                              `json:"username"`
		Email         string                              `json:"email"`
		Country       string                              `json:"country,omitempty"`
		City          string                              `json:"city,omitempty"`
		District      string                              `json:"district,omitempty"`
		Phone         string                              `json:"phone"`
		Whatsapp      string                              `json:"whatsapp"`
		CreatedAt     time.Time                           `json:"created_at"`
		AverageRating sqlc.GetAverageRatingForProviderRow `json:"average_rating"`
		AllRating     []sqlc.GetReviewsByProviderIDRow    `json:"all_rating"`
	}{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Country:       user.Country.String,
		City:          user.City.String,
		District:      user.District.String,
		Phone:         user.Phone,
		Whatsapp:      user.Whatsapp,
		CreatedAt:     user.CreatedAt,
		AverageRating: averageRating,
		AllRating: allRating,
	}

	ctx.JSON(http.StatusOK, userResponse)

}
