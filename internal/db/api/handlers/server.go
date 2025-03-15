package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/internal/db"
)

type Server struct {
	store  *db.Store // Взаимодействие с базой данных
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Настраиваем доверенные прокси
	router.SetTrustedProxies([]string{
		"127.0.0.1",      // локальный прокси
		"10.0.0.0/8",     // внутренняя сеть
		"172.16.0.0/12",  // Docker сети
		"192.168.0.0/16", // локальные сети
	})

	router.POST("/create-user", server.createUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}