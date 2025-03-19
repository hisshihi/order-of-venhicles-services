package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) loginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login", gin.H{
		"Title":      "Вход в систему",
		"ActivePage": "login",
	})
}
