package qr

import (
	"net/http"

	"github.com/isd-sgcu/rpkm66-gateway/app/dto"
	validate "github.com/isd-sgcu/rpkm66-gateway/app/validator"
	"github.com/isd-sgcu/rpkm66-gateway/pkg/rctx"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
)

type IService interface {
	CheckinVerify(string, int) (*proto.CheckinVerifyResponse, *dto.ResponseErr)
	CheckinConfirm(token string) (*proto.CheckinConfirmResponse, *dto.ResponseErr)
}

type Handler struct {
	checkinService IService
	validate       *validate.DtoValidator
}

// qr checkin which checkin for event day
// @Summary Get Token
// @Description get token by providing id return object has checkin type as enum (not to be confused with event type) 1 is checkin, 2 is checkout
// @Param event_type body dto.CheckinVerifyRequest true "event type (1 is Main event, 2 is Freshy Night)"
// @Tags QR
// @Accept json
// @Produce json
// @Success 200 {object} dto.CheckinVerifyResponse OK
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid body request
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 500 {object} dto.ResponseInternalErr Internal server error
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Router /qr/checkin/verify [post]
// @Security     AuthToken
func (h *Handler) CheckinVerify(ctx rctx.Context) {
	userid := ctx.UserID()
	cvr := &dto.CheckinVerifyRequest{}

	err := ctx.Bind(cvr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	res, errRes := h.checkinService.CheckinVerify(userid, cvr.EventType)

	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func NewHandler(checkinService IService, v *validate.DtoValidator) *Handler {
	return &Handler{
		checkinService: checkinService,
		validate:       v,
	}
}

// qr checkin which checkin for event day
// @Summary Confirm Checkin
// @Description Use token to confirm checkin
// @Param token body dto.CheckinConfirmRequest true "Token generated from CheckinVerify"
// @Tags QR
// @Accept json
// @Produce json
// @Success 200 {object} proto.CheckinConfirmResponse
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid body request
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 403 {object} dto.ResponseForbiddenErr Forbidden
// @Failure 500 {object} dto.ResponseInternalErr Internal server error
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Router /qr/checkin/confirm [post]
// @Security     AuthToken
func (h *Handler) CheckinConfirm(ctx rctx.Context) {
	ccr := &dto.CheckinConfirmRequest{}

	err := ctx.Bind(ccr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	res, errRes := h.checkinService.CheckinConfirm(ccr.Token)

	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
