package baan

import (
	"net/http"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/pkg/rctx"
	"github.com/isd-sgcu/rpkm66-gateway/pkg/service/baan"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
)

type Handler struct {
	service     baan.Service
	userService IUserService
}

type IUserService interface {
	FindOne(string) (*proto.User, *dto.ResponseErr)
}

func NewHandler(service baan.Service, userService IUserService) *Handler {
	return &Handler{service: service, userService: userService}
}

// FindAll is a function that get all baans
// @Summary Get all baans
// @Description Return the array of baan dto if successfully
// @Tags baan
// @Accept json
// @Produce json
// @Success 200 {object} []proto.Baan
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /baan [get]
func (h *Handler) FindAll(c rctx.Context) {
	result, err := h.service.FindAll()
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// FindOne is a function that get the baan data by id
// @Summary Get the baan data by id
// @Description Return the baan dto if successfully
// @Param id path string true "id"
// @Tags baan
// @Accept json
// @Produce json
// @Success 200 {object} proto.Baan
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid body request
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 404 {object} dto.ResponseNotfoundErr Not found baan
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /baan/{id} [get]
func (h *Handler) FindOne(c rctx.Context) {
	id, err := c.ID()
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ID",
		})

		return
	}

	result, errRes := h.service.FindOne(id)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserBaan is a function that get the baan data by id
// @Summary Get the user's baan
// @Description Return the baan dto if successfully
// @Tags baan
// @Accept json
// @Produce json
// @Success 200 {object} proto.Baan
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 404 {object} dto.ResponseNotfoundErr Not found baan
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /baan/user [get]
func (h *Handler) GetUserBaan(c rctx.Context) {
	id := c.UserID()

	usr, err := h.userService.FindOne(id)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	result, err := h.service.FindOne(usr.BaanId)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
