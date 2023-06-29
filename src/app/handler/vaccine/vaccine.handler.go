package vaccine

import (
	"net/http"

	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	validate "github.com/isd-sgcu/rpkm66-gateway/src/app/validator"
	"github.com/isd-sgcu/rpkm66-gateway/src/pkg/rctx"
)

type Handler struct {
	service  IService
	validate *validate.DtoValidator
}

type IService interface {
	Verify(string, string) (*dto.VaccineResponse, *dto.ResponseErr)
}

func NewHandler(service IService, validate *validate.DtoValidator) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}

// Verify is a function that verify the user status
// @Summary Verify the user status
// @Description Return nothing if success
// @Param user body dto.Verify true "verify dto"
// @Tags vaccine
// @Accept json
// @Produce json
// @Success 204 {bool} true
// @Failure 400 {object} dto.ResponseForbiddenErr Already verified
// @Failure 403 {object} dto.ResponseForbiddenErr Invalid phase
// @Security     AuthToken
// @Router /vaccine/verify [post]
func (h *Handler) Verify(ctx rctx.Context) {
	userId := ctx.UserID()

	verifyReq := dto.Verify{}
	err := ctx.Bind(&verifyReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input",
			Data:       nil,
		})
		return
	}

	if errors := h.validate.Validate(verifyReq); errors != nil {
		ctx.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid body request",
			Data:       errors,
		})
		return
	}

	ok, errRes := h.service.Verify(verifyReq.HCert, userId)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusNoContent, ok)
	return
}
