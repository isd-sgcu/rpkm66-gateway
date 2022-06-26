package auth

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	validate "github.com/isd-sgcu/rnkm65-gateway/src/app/validator"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"net/http"
)

type Handler struct {
	service    IService
	usrService IUserService
	validate   *validate.DtoValidator
}

type IService interface {
	VerifyTicket(string) (*proto.Credential, *dto.ResponseErr)
	Validate(string) (*dto.TokenPayloadAuth, *dto.ResponseErr)
	RefreshToken(string) (*proto.Credential, *dto.ResponseErr)
}

type IUserService interface {
	FindOne(string) (*proto.User, *dto.ResponseErr)
}

type IContext interface {
	Bind(interface{}) error
	JSON(int, interface{})
	UserID() string
}

func NewHandler(service IService, usrService IUserService, validate *validate.DtoValidator) *Handler {
	return &Handler{
		service:    service,
		usrService: usrService,
		validate:   validate,
	}
}

func (h *Handler) VerifyTicket(c IContext) {
	verifyTicket := dto.VerifyTicket{}

	err := c.Bind(&verifyTicket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error while binding the ticket: " + err.Error(),
			Data:       nil,
		})
		return
	}

	credential, errRes := h.service.VerifyTicket(verifyTicket.Ticket)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, credential)
}

func (h *Handler) Validate(c IContext) {
	userId := c.UserID()

	usr, err := h.usrService.FindOne(userId)
	if err != nil {
		switch err.StatusCode {
		case http.StatusNotFound:
			c.JSON(http.StatusUnauthorized, &dto.ResponseErr{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid user",
			})
		default:
			c.JSON(err.StatusCode, err)
		}
		return
	}

	c.JSON(http.StatusOK, usr)
}

func (h *Handler) RefreshToken(c IContext) {
	refreshToken := dto.RedeemNewToken{}

	err := c.Bind(&refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error while binding the ticket: " + err.Error(),
			Data:       nil,
		})
		return
	}

	credential, errRes := h.service.RefreshToken(refreshToken.RefreshToken)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, credential)
}
