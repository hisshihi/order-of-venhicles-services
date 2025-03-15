package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
)

type createUserParams struct {
	Username     string         `json:"username" binding:"required"`
	Email        string         `json:"email" binding:"required,email"`
	PasswordHash string         `json:"password_hash" binding:"min=6"`
	Country      sql.NullString `json:"country"`
	City         sql.NullString `json:"city"`
	District     sql.NullString `json:"district"`
	Phone        string         `json:"phone" binding:"required"`
	Whatsapp     string         `json:"whatsapp" binding:"required"`
}

// TODO: При создании добавить возвращение токена

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.CreateUserParams{
		Username: req.Username,
		Email: req.Email,
		PasswordHash: req.PasswordHash,
		Country: req.Country,
		City: req.City,
		District: req.District,
		Phone: req.Phone,
		Whatsapp: req.Whatsapp,
		Role: sqlc.NullRole{Role: sqlc.RoleClient, Valid: true},
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
