package file

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/constant"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	service     IService
	usrService  IUserService
	MaxFileSize int64
}

type IService interface {
	UploadImage(*dto.DecomposedFile) (string, *dto.ResponseErr)
}

type IUserService interface {
	FindOne(string) (*proto.User, *dto.ResponseErr)
}

type IContext interface {
	File(string, map[string]struct{}, int64) (*dto.DecomposedFile, error)
	JSON(int, interface{})
	UserID() string
}

func NewHandler(service IService, usrService IUserService, maxFileSize int) *Handler {
	return &Handler{
		service:     service,
		usrService:  usrService,
		MaxFileSize: int64(maxFileSize * 1024 * 1024),
	}
}

// UploadImage is a function that upload the image
// @Summary Upload the image
// @Description Return the filename if successfully
// @Tags file
// @Accept mpfd
// @Produce json
// @Success 201 {object} dto.FileResponse
// @Failure 400 {object} dto.ResponseBadRequestErr Invalid file
// @Failure 401 {object} dto.ResponseUnauthorizedErr Unauthorized
// @Failure 503 {object} dto.ResponseServiceDownErr Service is down
// @Failure 504 {object} dto.ResponseGatewayTimeoutErr Gateway timeout
// @Security     AuthToken
// @Router /file/image [post]
func (h *Handler) UploadImage(c IContext) {
	id := c.UserID()
	file, err := c.File(constant.Image, constant.AllowImageContentType, h.MaxFileSize)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "file").
			Str("module", "upload image").
			Msg("Invalid file")
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid file",
		})
		return
	}

	usr, errRes := h.usrService.FindOne(id)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	file.Filename = usr.StudentID

	filename, errRes := h.service.UploadImage(file)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusCreated, &dto.FileResponse{
		Filename: filename,
	})
}
