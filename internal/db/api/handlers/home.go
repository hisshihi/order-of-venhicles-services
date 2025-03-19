package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *Server) homePage(ctx *gin.Context) {
	categories, err := server.store.ListServiceCategories(ctx)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error", gin.H{
			"Title": "Ошибка",
			"Error": "Ошибка при получении категорий",
		})
		return
	}

	user, isAuthenticated := server.getUserDataFromTokenForHTML(ctx)

	ctx.HTML(http.StatusOK, "base", gin.H{
		"Title":           "Главная страница",
		"ActivePage":      "home",
		"Categories":      categories,
		"IsAuthenticated": isAuthenticated,
		"User":            user,
		"now":             time.Now(),
	})
}
