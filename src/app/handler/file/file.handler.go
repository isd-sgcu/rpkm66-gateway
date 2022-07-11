package file

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/constant/file"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type Handler struct {
	service     IService
	usrService  IUserService
	MaxFileSize int64
}

type IService interface {
	Upload(*dto.DecomposedFile, string, file.Tag, file.Type) (string, *dto.ResponseErr)
}

type IUserService interface {
	FindOne(string) (*proto.User, *dto.ResponseErr)
}

type IContext interface {
	File(string, map[string]struct{}, int64) (*dto.DecomposedFile, error)
	JSON(int, interface{})
	UserID() string
	GetFormData(string) string
}

func NewHandler(service IService, usrService IUserService, maxFileSize int) *Handler {
	return &Handler{
		service:     service,
		usrService:  usrService,
		MaxFileSize: int64(maxFileSize * 1024 * 1024),
	}
}

// Upload is a function that upload the image
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
func (h *Handler) Upload(c IContext) {
	id := c.UserID()

	tag := getTagNumber(c.GetFormData("tag"))
	if tag == file.UnknownTag {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid tag",
		})
		return
	}

	fileType := getTypeNumber(c.GetFormData("type"))
	if fileType == file.UnknownType {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid type",
		})
		return
	}

	content, err := c.File("file", file.AllowContentType, h.MaxFileSize)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "content").
			Str("module", "upload image").
			Msg("Invalid content")
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid content",
		})
		return
	}

	usr, errRes := h.usrService.FindOne(id)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	content.Filename = usr.StudentID

	filename, errRes := h.service.Upload(content, id, tag, fileType)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusCreated, &dto.FileResponse{
		Url: filename,
	})
}

func getTagNumber(tag string) file.Tag {
	switch strings.ToLower(tag) {
	case "profile":
		return file.Profile
	case "baan":
		return file.Baan
	default:
		return file.UnknownTag
	}
}

func getTypeNumber(fileType string) file.Type {
	switch strings.ToLower(fileType) {
	case "image":
		return file.Image
	case "file":
		return file.File
	default:
		return file.UnknownType
	}
}
