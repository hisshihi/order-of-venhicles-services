package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/hisshihi/order-of-venhicles-services/internal/config"
	"github.com/hisshihi/order-of-venhicles-services/internal/db"
	"github.com/hisshihi/order-of-venhicles-services/pkg/util"
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
	// Здесь можно добавить другие маршруты для администраторов
	// adminRoutes.GET("/users", server.listAllUsers)

	server.router = router
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

		fiels := strings.Fields(authorizationHeader)
		if len(fiels) < 2 {
			err := errors.New("неверный формат заголовка")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fiels[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("неподдерживаемый тип авторизации %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fiels[1]
		payload, err := server.maker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
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
