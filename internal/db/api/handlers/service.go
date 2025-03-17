package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
	"github.com/lib/pq"
)

type createServiceRequest struct {
	CategoryID  int64  `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description" binding:"required,min=50"`
	Price       string `json:"price" binding:"required"`
	Country string `json:"country" binding:"required"`
	City string `json:"city" binding:"required"`
	District string `json:"district" binding:"required"`
}

type createServiceResponse struct {
	ID          int64  `json:"id"`
	ProviderID  int64  `json:"provider_id"`
	CategoryID  int64  `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Country string `json:"country" binding:"required"`
	City string `json:"city" binding:"required"`
	District string `json:"district"`
	ProviderName string `json:"provider_name"`
	ProviderPhoto string `json:"provider_photo"`
	ProviderPhone string `json:"provider_phone"`
	ProviderWhatsapp string `json:"provider_whatsapp"`
	CategoryName string `json:"category_name"`
	ReviewsCount int64 `json:"reviews_count"`
	AverageRating float64 `json:"average_rating"`
}

func (server *Server) createService(ctx *gin.Context) {
	var req createServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.CreateServiceParams{
		ProviderID:  user.ID,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Country: sql.NullString{String: req.Country, Valid: true},
		City: sql.NullString{String: req.City, Valid: true},
		District: sql.NullString{String: req.District, Valid: true},
	}

	service, err := server.store.CreateService(ctx, arg)
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

	rsp := createServiceResponse{
		ID:          service.ID,
		ProviderID:  service.ProviderID,
		CategoryID:  service.CategoryID,
		Title:       service.Title,
		Description: service.Description,
		Price:       service.Price,
		Country: service.Country.String,
		City: service.City.String,
		District: service.District.String,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getServiceByProviderID(ctx *gin.Context) {
	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	service, err := server.store.GetServiceByID(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createServiceResponse{
		ID:          service.ID,
		ProviderID:  service.ProviderID,
		CategoryID:  service.CategoryID,
		Title:       service.Title,
		Description: service.Description,
		Price:       service.Price,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getServiceByIDRequest struct {
	ID int64 `uri:"id" binding:"min=1,required"`
}

func (server *Server) getServiceByID(ctx *gin.Context) {
	var req getServiceByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	service, err := server.store.GetServiceByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, service)
}

type listServicesRequest struct {
	ProviderID int64 `form:"provider_id" binding:"min=1,required"`
	PageID   int32 `form:"page_id" binding:"min=1,required"`
	PageSize int32 `form:"page_size" binding:"min=5,max=10,required"`
}

func (server *Server) listServiceByProviderID(ctx *gin.Context) {
	var req listServicesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.GetServicesByProviderIDParams{
		ProviderID: req.ProviderID,
		Limit:      int64(req.PageSize),
		Offset:     int64((req.PageID - 1) * req.PageSize),
	}

	services, err := server.store.GetServicesByProviderID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, services)
}

type listServicesByCategoryIDRequest struct {
	CategoryID int64 `form:"category_id" binding:"min=1,required"`
	PageID     int32 `form:"page_id" binding:"min=1,required"`
	PageSize   int32 `form:"page_size" binding:"min=5,max=10,required"`
}

func (server *Server) listServiceByCategoryID(ctx *gin.Context) {
	var req listServicesByCategoryIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.ListServicesByCategoryParams{
		CategoryID: req.CategoryID,
		Limit:      int64(req.PageSize),
		Offset:     int64((req.PageID - 1) * req.PageSize),
	}

	services, err := server.store.ListServicesByCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, services)
}

func (server *Server) listService(ctx *gin.Context) {
	var req listServicesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.ListServicesParams{
		Limit:  int64(req.PageSize),
		Offset: int64((req.PageID - 1) * req.PageSize),
	}

	services, err := server.store.ListServices(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, services)
}

type updateServiceRequest struct {
	CategoryID  int64  `json:"category_id" binding:"min=1,required"`
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description" binding:"required,min=50"`
	Price       string `json:"price" binding:"required"`
	Country     string `json:"country" binding:"required"`
	City        string `json:"city" binding:"required"`
	District    string `json:"district" binding:"required"`
}

func (server *Server) updateService(ctx *gin.Context) {
	var reqID getServiceByIDRequest
	var req updateServiceRequest

	if err := ctx.ShouldBindUri(&reqID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	service, err := server.store.GetServiceByID(ctx, reqID.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("услуга не найдена")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if service.ProviderID != user.ID {
		ctx.JSON(http.StatusForbidden,
			errorResponse(errors.New("у вас нет прав на обновление этой услуги")))
		return
	}

	categoryID, err := server.store.GetServiceCategoryByID(ctx, req.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("указанная категория не существует")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if _, err := strconv.ParseFloat(req.Price, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("некорректный формат цены")))
		return
	}

	arg := sqlc.UpdateServiceParams{
		ID:          service.ID,
		ProviderID:  service.ProviderID,
		CategoryID:  categoryID.ID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Country: sql.NullString{String: req.Country, Valid: true},
		City: sql.NullString{String: req.City, Valid: true},
		District: sql.NullString{String: req.District, Valid: true},
	}

	updatedService, err := server.store.UpdateService(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedService)
}

func (server *Server) deleteService(ctx *gin.Context) {
	var req getServiceByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserDataFromToken(ctx)
	if err != nil {
		return
	}

	service, err := server.store.GetServiceByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("услуга не найдена")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if service.ProviderID != user.ID {
		ctx.JSON(http.StatusForbidden,
			errorResponse(errors.New("у вас нет прав на удаление этой услуги")))
		return
	}

	arg := sqlc.DeleteServiceParams{
		ID: service.ID,
		ProviderID: service.ProviderID,
	}

	err = server.store.DeleteService(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusNotFound, nil)
}
