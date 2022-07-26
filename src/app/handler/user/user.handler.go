package user

import (
	"fmt"

	"net/http"

	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/estamp"
	validate "github.com/isd-sgcu/rnkm65-gateway/src/app/validator"
	"github.com/isd-sgcu/rnkm65-gateway/src/interfaces/qr"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
)

type Handler struct {
	service  IService
	validate *validate.DtoValidator
}

type IContext interface {
	JSON(int, interface{})
	UserID() string
	Bind(interface{}) error
	ID() (string, error)
}

func NewHandler(service IService, validate *validate.DtoValidator) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}

type IService interface {
	FindOne(string) (*proto.User, *dto.ResponseErr)
	Create(*dto.UserDto) (*proto.User, *dto.ResponseErr)
	Update(string, *dto.UpdateUserDto) (*proto.User, *dto.ResponseErr)
	CreateOrUpdate(*dto.UserDto) (*proto.User, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
	GetUserEstamp(string) (*proto.GetUserEstampResponse, *dto.ResponseErr)
	ConfirmEstamp(string, string) (*proto.ConfirmEstampResponse, *dto.ResponseErr)
}

// FindOne is a function that get the user data by id
// @Summary Get the user data by id
// @Description Return the user dto if successfully
// @Param id path string true "id"
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} proto.User
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid body request
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 404 {object} dto.ResponseNotfoundErr Not found user
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Router /user/{id} [get]
func (h *Handler) FindOne(ctx IContext) {
	id, err := ctx.ID()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid ID",
			Data:       nil,
		})
		return
	}

	user, errRes := h.service.FindOne(id)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, user)
	return
}

// Create is a function that create new user
// @Summary Create new user
// @Description Return the user dto if successfully
// @Param user body dto.UserDto true "User DTO"
// @Tags user
// @Accept json
// @Produce json
// @Success 201 {object} proto.User
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid request body
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 403 {object} dto.ResponseForbiddenErr Insufficiency permission to create user
// @Failure 404 {object} dto.ResponseNotfoundErr Not found user
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /user [post]
func (h *Handler) Create(ctx IContext) {
	usrDto := dto.UserDto{}

	err := ctx.Bind(&usrDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if errors := h.validate.Validate(usrDto); errors != nil {
		ctx.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid body request",
			Data:       errors,
		})
		return
	}

	user, errRes := h.service.Create(&usrDto)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusCreated, user)
	return
}

// Update is a function that Update the user data if exists
// @Summary Update the user data if exists
// @Description Return the user dto if successfully
// @Param user body dto.UserDto true "user dto"
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} proto.User
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid request body
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /user [patch]
func (h *Handler) Update(ctx IContext) {
	usrId := ctx.UserID()

	usrDto := dto.UpdateUserDto{}

	err := ctx.Bind(&usrDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	user, errRes := h.service.Update(usrId, &usrDto)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, user)
	return
}

// CreateOrUpdate is a function that Create new user if it doesn't exist and Update the user data if exists
// @Summary Create new user if it doesn't exist and Update the user data if exists
// @Description Return the user dto if successfully
// @Param user body dto.UpdateUserDto true "user dto"
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} proto.User
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid request body
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /user [put]
func (h *Handler) CreateOrUpdate(ctx IContext) {
	id := ctx.UserID()
	usrDto := dto.UserDto{}

	err := ctx.Bind(&usrDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	usrDto.ID = id

	if errors := h.validate.Validate(usrDto); errors != nil {
		for _, response := range errors {
			fmt.Println(response)
		}

		ctx.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid body request",
			Data:       errors,
		})
		return
	}

	user, errRes := h.service.CreateOrUpdate(&usrDto)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, user)
	return
}

// Delete is a function that delete the user
// @Summary Delete the user
// @Description Return the user dto if successfully
// @Param id path string true "id"
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {bool} true
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid ID
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 403 {object} dto.ResponseForbiddenErr Insufficiency permission to delete user
// @Failure 404 {object} dto.ResponseNotfoundErr Not found user
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /user/{id} [delete]
func (h *Handler) Delete(ctx IContext) {
	id, err := ctx.ID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	user, errRes := h.service.Delete(id)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, user)
	return
}

// Get estamp overview id on what user has
// @Summary Get user estamp
// @Description Get estamp id overview on what user has *Return {} with success status code if user has no estamp
// @Tags event
// @Accept json
// @Produce json
// @Success 200 {object} proto.GetUserEstampResponse OK
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 500 {object} dto.ResponseInternalErr Internal server error
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Router /estamp/user [get]
// @Security     AuthToken
func (h *Handler) GetUserEstamp(ctx estamp.IContext) {
	id := ctx.UserID()

	res, errRes := h.service.GetUserEstamp(id)

	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
	return
}

// get estamp for user
// @Summary Confirm Estamp
// @Description get estamp
// @Param token body dto.ConfirmEstampRequest true "Event id"
// @Tags QR
// @Accept json
// @Produce json
// @Success 204 {object} proto.ConfirmEstampResponse OK
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid body request
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 500 {object} dto.ResponseInternalErr Internal server error
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Router /qr/estamp/confirm [post]
// @Security     AuthToken
func (h *Handler) ConfirmEstamp(ctx qr.IContext) {
	userid := ctx.UserID()
	ce := &dto.ConfirmEstampRequest{}

	err := ctx.Bind(ce)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	res, errRes := h.service.ConfirmEstamp(userid, ce.EventId)

	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusNoContent, res)
	return
}
