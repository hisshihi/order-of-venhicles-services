package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/hisshihi/order-of-venhicles-services/internal/config"
	"github.com/hisshihi/order-of-venhicles-services/internal/db"
	"github.com/hisshihi/order-of-venhicles-services/pkg/util"
	"github.com/lib/pq"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type Server struct {
	config config.Config
	store  *db.Store // Взаимодействие с базой данных
	router *gin.Engine
	maker  util.Maker // Единственное поле для работы с токенами
}

func NewServer(config config.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := util.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store:  store,
		maker:  tokenMaker,
	}

	server.setupServer()
	return server, nil
}

func (server *Server) setupServer() {
	router := gin.Default()

	// Настраиваем доверенные прокси
	router.SetTrustedProxies([]string{
		"127.0.0.1",      // локальный прокси
		"10.0.0.0/8",     // внутренняя сеть
		"172.16.0.0/12",  // Docker сети
		"192.168.0.0/16", // локальные сети
	})

	// Публичные маршруты (без авторизации)
	router.POST("/create-user", server.createUser)
	router.POST("/user/login", server.loginUser)

	// Защищённые маршруты с ролевым доступом
	authRoutes := router.Group("/")
	authRoutes.Use(server.authMiddleware())

	// Маршруты доступные всем авторизированным пользователям
	authRoutes.GET("/users/me", server.getCurrentUser)

	// Маршруты для клиентов
	clientRoutes := router.Group("/client")
	clientRoutes.Use(server.authMiddleware())
	clientRoutes.Use(server.roleCheckMiddleware(
		string(sqlc.RoleClient),
		string(sqlc.RoleProvider),
		string(sqlc.RolePartner),
		string(sqlc.RoleAdmin),
	))
	// Добавьте здесь маршруты для клиентов
	// clientRoutes.POST("/orders", server.createOrder)

	// Маршруты для поставщиков услуг
	providerRoutes := router.Group("/provider")
	providerRoutes.Use(server.authMiddleware())
	providerRoutes.Use(server.roleCheckMiddleware(
		string(sqlc.RoleProvider),
		string(sqlc.RoleAdmin),
	))
	// Добавьте здесь маршруты для поставщиков
	// providerRoutes.POST("/services", server.createService)

	// Маршруты для партнеров
	partnerRoutes := router.Group("/partner")
	partnerRoutes.Use(server.authMiddleware())
	partnerRoutes.Use(server.roleCheckMiddleware(
		string(sqlc.RolePartner),
		string(sqlc.RoleAdmin),
	))
	// Добавьте здесь маршруты для партнеров
	// partnerRoutes.POST("/promo-codes", server.createPromoCode)

	// Маршруты только для администраторов
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(server.authMiddleware())
	adminRoutes.Use(server.roleCheckMiddleware(string(sqlc.RoleAdmin)))
	adminRoutes.GET("/users/:id", server.getUserByID) // Доступ к данным пользователя по ID
	adminRoutes.POST("/category", server.createCategory)
	adminRoutes.GET("/category/:id", server.getCategoryByID)
	adminRoutes.GET("/category", server.listCategory)
	adminRoutes.PUT("/category/:id", server.updateCategory)
	adminRoutes.DELETE("/category/:id", server.deleteCategory)
	// Здесь можно добавить другие маршруты для администраторов
	// adminRoutes.GET("/users", server.listAllUsers)

	server.router = router

	// Создаем администратора по умолчанию при инициализации сервера
	if err := server.createAdminDefault(); err != nil {
		// Логируем ошибку, но не прерываем инициализацию сервера
		fmt.Printf("Не удалось создать администратора по умолчанию: %v\n", err)
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// authMiddleware Middleware аунтификации (проверка токена)
func (server *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("заголовок авторизации не предоставлен")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("неверный формат заголовка")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("неподдерживаемый тип авторизации %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := server.maker.VerifyToken(accessToken)
		if err != nil {
			// Более информативные сообщения в зависимости от типа ошибки
			var errorMsg string
			if errors.Is(err, util.ErrExpiredToken) {
				errorMsg = "срок действия токена истек, необходимо пройти авторизацию повторно"
			} else if errors.Is(err, util.ErrInvalidToken) {
				errorMsg = "недействительный токен"
			} else {
				errorMsg = fmt.Sprintf("ошибка проверки токена: %v", err)
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New(errorMsg)))
			return
		}

		// Двойная проверка срока действия токена на уровне middleware
		if err := payload.Valid(); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				errorResponse(errors.New("срок действия токена истек, необходимо пройти авторизацию повторно")))
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}

// roleCheckMiddleware Middleware для првоерки роли
// Функция-генератор middleware, принимает список разрешённый ролей
func (server *Server) roleCheckMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем payload из контекста (установленный authMiddleware)
		payload, exists := c.Get(authorizationPayloadKey)
		if !exists {
			err := errors.New("требуется авторизация")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		tokenPayload, ok := payload.(*util.Payload)
		if !ok {
			err := errors.New("неверный тип payload")
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Проверяем, есть ли роль пользователя в списке разрешённых
		roleAllowed := false
		for _, role := range allowedRoles {
			if tokenPayload.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			err := errors.New("доступ запрешён: недостаточно прав")
			c.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}

		c.Next()
	}
}

type createAdminDefaultParam struct {
	Username     string  `json:"username"`
	Email        string  `json:"email" `
	PasswordHash string  `json:"password_hash"`
	Country      *string `json:"country"`
	City         *string `json:"city"`
	District     *string `json:"district"`
	Phone        string  `json:"phone"`
	Whatsapp     string  `json:"whatsapp"`
	Role         *string `json:"role"`
}

// Создание администратора по умолчанию
func (server *Server) createAdminDefault() error {
	hashedPassword, err := util.HashPassword(server.config.AdminPassword)
	if err != nil {
		return fmt.Errorf("ошибка при хэшировании пароля: %w", err)
	}

	arg := sqlc.CreateUserParams{
		Username:     server.config.AdminUsername,
		Email:        server.config.AdminEmail,
		PasswordHash: hashedPassword,
		Country:      sql.NullString{},
		City:         sql.NullString{},
		District:     sql.NullString{},
		Phone:        "+7-(999)-999-99-99",
		Whatsapp:     "+7-(999)-999-99-99",
		Role:         sqlc.NullRole{Role: sqlc.RoleAdmin, Valid: true},
	}

	// Создаем администратора в базе данных
	_, err = server.store.CreateUser(context.Background(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				// Администратор уже существует - это не ошибка при инициализации
				fmt.Println("Администратор уже существует в системе")
				return nil
			}
		}
		return fmt.Errorf("ошибка при создании администратора: %w", err)
	}

	fmt.Println("--------------------------")
	fmt.Println("Администратор успешно создан")
	fmt.Println("--------------------------")
	return nil
}
