package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/hisshihi/order-of-venhicles-services/internal/config"
	"github.com/hisshihi/order-of-venhicles-services/internal/db"
	"github.com/hisshihi/order-of-venhicles-services/pkg/util"
	"github.com/lib/pq"
	"golang.org/x/time/rate"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

var limiter = rate.NewLimiter(1, 300)

func rateLimiter(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, errorResponse(errors.New("too many requests")))
		c.Abort()
		return
	}
	c.Next()
}

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

	// Настройка CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// Применяем middleware для ограничения запросов
	router.Use(rateLimiter)

	// Публичные маршруты (без авторизации)
	router.POST("/create-user", server.createUser)
	router.POST("/user/login", server.loginUser)
	router.GET("/categories", server.listCategory)
	router.GET("/categories/:id", server.getCategoryByID)

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
	clientRoutes.POST("/orders", server.createOrder)
	clientRoutes.GET("/orders/:id", server.getOrderByID)
	clientRoutes.POST("/orders/:id/status", server.updateOrderStatus)
	clientRoutes.GET("/orders/", server.listOrders)
	clientRoutes.POST("/reviews", server.createReview)
	clientRoutes.GET("/reviews", server.listReviewByProviderID)
	clientRoutes.GET("/reviews/:id/rating", server.getAverageRatingForProvider)
	clientRoutes.DELETE("/reviews/:id", server.deleteReview)

	clientRoutes.GET("/services/:id", server.getServiceByID)
	clientRoutes.GET("/services/list", server.listService)
	clientRoutes.GET("/services/list/provider", server.listServiceByProviderID)
	clientRoutes.GET("/services/list/category", server.listServiceByCategoryID)

	// Маршруты для поставщиков услуг
	providerRoutes := router.Group("/provider")
	providerRoutes.Use(server.authMiddleware())
	providerRoutes.Use(server.roleCheckMiddleware(
		string(sqlc.RoleProvider),
		string(sqlc.RoleAdmin),
	))

	// Маршруты, которые НЕ требуют подписки
	// TODO: убрать конечные точки создание и обновление и реализовать эти методы только после оплаты(статус 200)
	// providerRoutes.POST("/subscriptions", server.createSubscription)
	// providerRoutes.POST("/subscriptions/update", server.updateSubscription)
	// providerRoutes.GET("/subscriptions/check", server.checkSubscriptionActive)
	// providerRoutes.POST("/payment", server.createPayment)
	providerRoutes.POST("/payments/initiate", server.initiateSubscriptionPayment)
	providerRoutes.POST("/payments/callback", server.processPaymentCallback)
	providerRoutes.GET("/payments/:payment_id/status", server.checkPaymentStatus)
	// TODO: удалить после тестов
	providerRoutes.POST("/payments/simulate", server.simuldateSuccessfulPayment)

	// Маршруты, которые требуют подписку
	subscriptionRequiredRoutes := providerRoutes.Group("/")
	subscriptionRequiredRoutes.Use(server.subscriptionCheckMiddleware())

	subscriptionRequiredRoutes.POST("/services", server.createService)
	subscriptionRequiredRoutes.GET("/services", server.getServiceByProviderID)
	subscriptionRequiredRoutes.GET("/services/list/u", server.listServiceFromProvider)
	subscriptionRequiredRoutes.PUT("/services/:id", server.updateService)
	subscriptionRequiredRoutes.DELETE("/services/:id", server.deleteService)
	subscriptionRequiredRoutes.POST("/orders/:id/accept", server.acceptOrder)
	subscriptionRequiredRoutes.GET("/orders/available", server.listAvailableOrders)
	subscriptionRequiredRoutes.GET("/orders/statistics", server.getOrdersStatistics)
	subscriptionRequiredRoutes.GET("/orders/category/:category_id", server.getOrdersByCategory)
	subscriptionRequiredRoutes.GET("/reviews/only", server.getReviewsByThisProviderID)

	// Маршруты для партнеров
	partnerRoutes := router.Group("/partner")
	partnerRoutes.Use(server.authMiddleware())
	partnerRoutes.Use(server.roleCheckMiddleware(
		string(sqlc.RolePartner),
		string(sqlc.RoleAdmin),
	))
	// Добавьте здесь маршруты для партнеров
	partnerRoutes.POST("/promo-codes", server.createPromoCode)
	partnerRoutes.GET("/subscriptions/provider", server.listSubsciptionsByProviderID)

	// Маршруты только для администраторов
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(server.authMiddleware())
	adminRoutes.Use(server.roleCheckMiddleware(string(sqlc.RoleAdmin)))
	adminRoutes.GET("/users/:id", server.getUserByID) // Доступ к данным пользователя по ID
	adminRoutes.POST("/category", server.createCategory)
	adminRoutes.PUT("/category/:id", server.updateCategory)
	adminRoutes.DELETE("/category/:id", server.deleteCategory)
	// Здесь можно добавить другие маршруты для администраторов
	// adminRoutes.GET("/users", server.listAllUsers)

	server.router = router

	// Создаем пользователей по умолчанию при инициализации приложения
	if err := server.createDefaultUsers(); err != nil {
		// Логируем ошибку, но не прерываем инициализацию сервера
		fmt.Printf("Не удалось создать пользователей по умолчанию: %v\n", err)
	}
}

func (server *Server) Start(address string) error {
	server.startSubscriptionChecker()

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

// Создание пользователей по умолчанию при инициализации приложения
func (server *Server) createDefaultUsers() error {
	// 1. Создание администратора
	if err := server.createUserWithRole(
		server.config.AdminUsername,
		server.config.AdminEmail,
		server.config.AdminPassword,
		"+7-(999)-999-99-99",
		"+7-(999)-999-99-99",
		sqlc.RoleAdmin,
	); err != nil {
		return fmt.Errorf("ошибка при создании администратора: %w", err)
	}

	// 2. Создание клиента
	if err := server.createUserWithRole(
		server.config.ClientUsername,
		server.config.ClientEmail,
		server.config.ClientPassword,
		"+7-(888)-888-88-88",
		"+7-(888)-888-88-88",
		sqlc.RoleClient,
	); err != nil {
		return fmt.Errorf("ошибка при создании клиента: %w", err)
	}

	// 3. Создание провайдера услуг
	if err := server.createUserWithRole(
		server.config.ProviderUsername,
		server.config.ProviderEmail,
		server.config.ProviderPassword,
		"+7-(777)-777-77-77",
		"+7-(777)-777-77-77",
		sqlc.RoleProvider,
	); err != nil {
		return fmt.Errorf("ошибка при создании провайдера: %w", err)
	}

	// 4. Создание партнера
	if err := server.createUserWithRole(
		server.config.PartnerUsername,
		server.config.PartnerEmail,
		server.config.PartnerPassword,
		"+7-(666)-666-66-66",
		"+7-(666)-666-66-66",
		sqlc.RolePartner,
	); err != nil {
		return fmt.Errorf("ошибка при создании партнера: %w", err)
	}

	return nil
}

// Вспомогательная функция для создания пользователя с определенной ролью
func (server *Server) createUserWithRole(username, email, password, phone, whatsapp string, role sqlc.Role) error {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return fmt.Errorf("ошибка при хэшировании пароля: %w", err)
	}

	arg := sqlc.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Country:      sql.NullString{},
		City:         sql.NullString{},
		District:     sql.NullString{},
		Phone:        phone,
		Whatsapp:     whatsapp,
		Role:         sqlc.NullRole{Role: role, Valid: true},
	}

	// Создаем пользователя в базе данных
	_, err = server.store.CreateUser(context.Background(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				// Пользователь уже существует - это не ошибка при инициализации
				fmt.Printf("Пользователь %s (%s) уже существует в системе\n", username, role)
				return nil
			}
		}
		return err
	}

	fmt.Printf("Пользователь %s с ролью %s успешно создан\n", username, role)
	return nil
}

func (server *Server) getUserDataFromToken(ctx *gin.Context) (sqlc.User, error) {
	// Получаем payload из токена авторизации
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		err := errors.New("требуется авторизация")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return sqlc.User{}, err
	}

	// Приводим payload к нужному типу
	tokenPayload, ok := payload.(*util.Payload)
	if !ok {
		err := errors.New("неверный тип payload")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return sqlc.User{}, err
	}

	// Получаем пользователя по ID из токена
	user, err := server.store.GetUserByIDFromUser(ctx, tokenPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return sqlc.User{}, err
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return sqlc.User{}, err
	}

	return user, nil
}

func (server *Server) updateExpiredSubscriptions(ctx *gin.Context) error {
	// Вызываем SQL-запрос для проверки истёкших подписок
	expiredSubs, err := server.store.CheckAndUpdateExpiredSubscriptions(ctx)
	if err != nil {
		log.Printf("Ошибка при обновлении истёкших подписок %v", err)
		return err
	}

	log.Printf("Обновлено %d истёкших подписок", len(expiredSubs))

	for _, sub := range expiredSubs {
		log.Printf("Подписка ID: %d для провайдера ID: %d истекла %s", sub.ID, sub.ProviderID, sub.EndDate.Format("2006-01-02"))

		// TODO: добавить отправку уведомления
	}

	return nil
}
