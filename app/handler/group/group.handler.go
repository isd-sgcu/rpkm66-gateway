package group

import (
	"net/http"
	"net/url"

	"github.com/isd-sgcu/rpkm66-gateway/app/dto"
	validate "github.com/isd-sgcu/rpkm66-gateway/app/validator"
	"github.com/isd-sgcu/rpkm66-gateway/pkg/rctx"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
)

type Handler struct {
	service  IService
	validate *validate.DtoValidator
}

func NewHandler(service IService, validate *validate.DtoValidator) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}

type IService interface {
	FindOne(string) (*proto.Group, *dto.ResponseErr)
	FindByToken(string) (*proto.FindByTokenGroupResponse, *dto.ResponseErr)
	Update(*dto.GroupDto, string) (*proto.Group, *dto.ResponseErr)
	Join(string, string) (*proto.Group, *dto.ResponseErr)
	DeleteMember(string, string) (*proto.Group, *dto.ResponseErr)
	Leave(string) (*proto.Group, *dto.ResponseErr)
	SelectBaan(string, []string) (bool, *dto.ResponseErr)
}

// FindOne is a function that get the group data
// @Summary Get the group data
// @Description Return the group dto if successfully
// @Tags group
// @Accept json
// @Produce json
// @Success 200 {object} proto.Group
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 404 {object} dto.ResponseNotfoundErr Not found group
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /group [get]
func (h *Handler) FindOne(ctx rctx.Context) {
	userId := ctx.UserID()

	group, errRes := h.service.FindOne(userId)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, group)
	return
}

// FindByToken is a function that get the group data by token
// @Summary Get the group data by token
// @Description Return the group dto if successfully
// @Param token path string true "token"
// @Tags group
// @Accept json
// @Produce json
// @Success 200 {object} proto.FindByTokenGroupResponse
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 404 {object} dto.ResponseNotfoundErr Not found group
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /group/{token} [get]
func (h *Handler) FindByToken(ctx rctx.Context) {
	tokenUrl := ctx.Param("token")

	token, _ := url.QueryUnescape(tokenUrl)

	res, errRes := h.service.FindByToken(token)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
	return
}

// Join is a function for user to join the group
// @Summary Join the existing group
// @Description Return the group dto if successfully
// @Param token path string true "token"
// @Tags group
// @Accept json
// @Produce json
// @Success 200 {object} proto.Group
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid ID or Request Body
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 403 {object} dto.ResponseForbiddenErr Insufficiency permission to join group
// @Failure 404 {object} dto.ResponseNotfoundErr Not found group
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /group/{token} [post]
func (h *Handler) Join(ctx rctx.Context) {
	tokenUrl := ctx.Param("token")

	userId := ctx.UserID()
	token, _ := url.QueryUnescape(tokenUrl)
	group, errRes := h.service.Join(token, userId)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, group)
	return

}

// DeleteMember is a function that delete member from the group
// @Summary Delete member from the group
// @Description Return the group dto if successfully
// @Param member_id path string true "member_id"
// @Tags group
// @Accept json
// @Produce json
// @Success 200 {object} proto.Group
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid ID
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 403 {object} dto.ResponseForbiddenErr Insufficiency permission to delete group
// @Failure 404 {object} dto.ResponseNotfoundErr Not found group
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /group/members/{member_id} [delete]
func (h *Handler) DeleteMember(ctx rctx.Context) {
	userId := ctx.Param("member_id")

	leaderId := ctx.UserID()

	group, errRes := h.service.DeleteMember(userId, leaderId)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, group)
	return
}

// Leave is a function for to leave the group
// @Summary Leave the current group and Create a new group
// @Description Return the group dto if successfully
// @Tags group
// @Accept json
// @Produce json
// @Success 200 {object} proto.Group
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid ID
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 403 {object} dto.ResponseForbiddenErr Insufficiency permission to leave group
// @Failure 404 {object} dto.ResponseNotfoundErr Not found group
// @Failure 500 {object} dto.ResponseInternalErr Internal error
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /group/leave [delete]
func (h *Handler) Leave(ctx rctx.Context) {
	userId := ctx.UserID()

	group, errRes := h.service.Leave(userId)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, group)
	return
}

// SelectBaan is a function for to select the baan
// @Summary select baan for the group (leader only)
// @Description Return nothing if successfully
// @Param selectBaanDto body dto.SelectBaan true "Select baan dto"
// @Tags group
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid ID
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 403 {object} dto.ResponseForbiddenErr Insufficiency permission to select the baan
// @Failure 404 {object} dto.ResponseNotfoundErr Not found group
// @Failure 500 {object} dto.ResponseInternalErr Internal error
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Security     AuthToken
// @Router /group/select [put]
func (h *Handler) SelectBaan(ctx rctx.Context) {
	userId := ctx.UserID()

	baans := dto.SelectBaan{}
	err := ctx.Bind(&baans)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid baan ids")
		return
	}

	success, errRes := h.service.SelectBaan(userId, baans.Baans)
	if !success {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusNoContent, struct{}{})
}
