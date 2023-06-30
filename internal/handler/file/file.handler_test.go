package file

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/constant/file"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	mock "github.com/isd-sgcu/rpkm66-gateway/mocks/file"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/rctx"
	mockUsr "github.com/isd-sgcu/rpkm66-gateway/mocks/user"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/user/v1"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTest struct {
	suite.Suite
	filename          string
	file              []byte
	key               string
	maxFileSize       int
	user              *proto.User
	fileDecompose     *dto.DecomposedFile
	GatewayTimeoutErr *dto.ResponseErr
	ServiceDownErr    *dto.ResponseErr
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(UserHandlerTest))
}

func (t *UserHandlerTest) SetupTest() {
	t.user = &proto.User{
		Id:              faker.UUIDDigit(),
		Title:           faker.Word(),
		Firstname:       faker.Word(),
		Lastname:        faker.Word(),
		Nickname:        faker.Word(),
		StudentID:       faker.Word(),
		Faculty:         faker.Word(),
		Year:            faker.Word(),
		Phone:           faker.Phonenumber(),
		LineID:          faker.Word(),
		Email:           faker.Email(),
		AllergyFood:     faker.Word(),
		FoodRestriction: faker.Word(),
		AllergyMedicine: faker.Word(),
		Disease:         faker.Word(),
		ImageUrl:        faker.URL(),
		CanSelectBaan:   false,
	}
	t.filename = faker.Word()
	t.file = []byte("Hello")

	t.key = "file"

	t.maxFileSize = 10

	t.fileDecompose = &dto.DecomposedFile{
		Filename: t.filename,
		Data:     t.file,
	}

	t.GatewayTimeoutErr = &dto.ResponseErr{
		StatusCode: http.StatusGatewayTimeout,
		Message:    "Connection timeout",
		Data:       nil,
	}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}
}

func (t *UserHandlerTest) TestUploadSuccess() {
	want := &dto.FileResponse{Url: t.filename}

	c := rctx.ContextMock{}
	c.On("File", t.key, file.AllowContentType).Return(t.fileDecompose, nil)
	c.On("UserID").Return(t.user.Id)
	c.On("GetFormData", "tag").Return("profile", nil)
	c.On("GetFormData", "type").Return("image", nil)

	usrSrv := mockUsr.ServiceMock{}
	usrSrv.On("FindOne", t.user.Id).Return(t.user, nil)

	srv := mock.ServiceMock{}
	srv.On("Upload", t.fileDecompose, t.user.Id, file.Tag(file.Profile), file.Type(file.Image)).Return(t.filename, nil)

	hdr := NewHandler(&srv, &usrSrv, t.maxFileSize)

	hdr.Upload(&c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUploadInvalidFile() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid content",
		Data:       nil,
	}

	c := rctx.ContextMock{}
	c.On("File", t.key, file.AllowContentType).Return(nil, errors.New("Invalid content"))
	c.On("UserID").Return(t.user.Id)
	c.On("GetFormData", "tag").Return("profile", nil)
	c.On("GetFormData", "type").Return("image", nil)

	usrSrv := mockUsr.ServiceMock{}
	usrSrv.On("FindOne", t.user.Id).Return(t.user, nil)

	srv := mock.ServiceMock{}
	srv.On("Upload", t.fileDecompose, t.user.Id, file.Profile, file.Type(file.Image)).Return("", nil)

	hdr := NewHandler(&srv, &usrSrv, t.maxFileSize)

	hdr.Upload(&c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUploadFailed() {
	testUploadFailed(t.T(), t.GatewayTimeoutErr, t.key, t.fileDecompose, t.maxFileSize, t.user)
	testUploadFailed(t.T(), t.ServiceDownErr, t.key, t.fileDecompose, t.maxFileSize, t.user)
}

func testUploadFailed(t *testing.T, err *dto.ResponseErr, key string, decomposedFile *dto.DecomposedFile, maxFileSize int, user *proto.User) {
	want := err

	c := rctx.ContextMock{}
	c.On("File", key, file.AllowContentType).Return(decomposedFile, nil)
	c.On("UserID").Return(user.Id)
	c.On("GetFormData", "tag").Return("profile", nil)
	c.On("GetFormData", "type").Return("image", nil)

	usrSrv := mockUsr.ServiceMock{}
	usrSrv.On("FindOne", user.Id).Return(user, nil)

	srv := mock.ServiceMock{}
	srv.On("Upload", decomposedFile, user.Id, file.Tag(file.Profile), file.Type(file.Image)).Return("", err)

	hdr := NewHandler(&srv, &usrSrv, maxFileSize)

	hdr.Upload(&c)

	assert.Equal(t, want, c.V)
}

func (t *UserHandlerTest) TestUploadGrpcErr() {
	want := t.ServiceDownErr

	c := rctx.ContextMock{}
	c.On("File", t.key, file.AllowContentType).Return(t.fileDecompose, nil)
	c.On("UserID").Return(t.user.Id)
	c.On("GetFormData", "tag").Return("profile", nil)
	c.On("GetFormData", "type").Return("image", nil)

	usrSrv := mockUsr.ServiceMock{}
	usrSrv.On("FindOne", t.user.Id).Return(nil, t.ServiceDownErr)

	srv := mock.ServiceMock{}
	srv.On("Upload", t.fileDecompose, t.user.Id, file.Profile, file.Type(file.Image)).Return("", nil)

	hdr := NewHandler(&srv, &usrSrv, t.maxFileSize)

	hdr.Upload(&c)

	assert.Equal(t.T(), want, c.V)
}
