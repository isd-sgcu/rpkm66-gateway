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

// VerifyTicket is a function that send ticket to verify at chula sso and generate the new credential
// @Summary Verify ticket and get credential
// @Description Return the credential if successfully
// @Param register body dto.VerifyTicket true "refresh token dto"
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.Credential
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Security     AuthToken
// @Router /auth/verify [post]
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

// Validate is a function check the user token and return user dto
// @Summary Check user status and user info
// @Description Return the user dto if successfully
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} proto.User
// @Failure 401 {object} dto.ResponseUnauthorizedErr "Invalid token"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Security     AuthToken
// @Router /auth/me [get]
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

// RefreshToken is a function that redeem new credentials
// @Summary Redeem new token
// @Description Return the credentials if successfully
// @Param register body dto.RedeemNewToken true "refresh token dto"
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.Credential
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 401 {object} dto.ResponseUnauthorizedErr "Invalid refresh token"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /auth/refreshToken [post]
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
