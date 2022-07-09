package file

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/constant"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/file"
	mockUsr "github.com/isd-sgcu/rnkm65-gateway/src/mocks/user"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type UserHandlerTest struct {
	suite.Suite
	filename          string
	file              []byte
	imgKey            string
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

	t.imgKey = constant.Image

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

func (t *UserHandlerTest) TestUploadImageSuccess() {
	want := &dto.FileResponse{Filename: t.filename}

	c := mock.ContextMock{}
	c.On("File", t.imgKey, constant.AllowImageContentType).Return(t.fileDecompose, nil)
	c.On("UserID").Return(t.user.Id)

	usrSrv := mockUsr.ServiceMock{}
	usrSrv.On("FindOne", t.user.Id).Return(t.user, nil)

	srv := mock.ServiceMock{}
	srv.On("UploadImage", t.fileDecompose).Return(t.filename, nil)

	hdr := NewHandler(&srv, &usrSrv, t.maxFileSize)

	hdr.UploadImage(&c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUploadImageInvalidFile() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid file",
		Data:       nil,
	}

	c := mock.ContextMock{}
	c.On("File", t.imgKey, constant.AllowImageContentType).Return(nil, errors.New("Invalid file"))
	c.On("UserID").Return(t.user.Id)

	usrSrv := mockUsr.ServiceMock{}
	usrSrv.On("FindOne", t.user.Id).Return(t.user, nil)

	srv := mock.ServiceMock{}
	srv.On("UploadImage", t.fileDecompose).Return("", nil)

	hdr := NewHandler(&srv, &usrSrv, t.maxFileSize)

	hdr.UploadImage(&c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUploadImageFailed() {
	testUploadFailed(t.T(), t.GatewayTimeoutErr, t.imgKey, t.fileDecompose, t.maxFileSize, t.user)
	testUploadFailed(t.T(), t.ServiceDownErr, t.imgKey, t.fileDecompose, t.maxFileSize, t.user)
}

func testUploadFailed(t *testing.T, err *dto.ResponseErr, key string, file *dto.DecomposedFile, maxFileSize int, user *proto.User) {
	want := err

	c := mock.ContextMock{}
	c.On("File", key, constant.AllowImageContentType).Return(file, nil)
	c.On("UserID").Return(user.Id)

	usrSrv := mockUsr.ServiceMock{}
	usrSrv.On("FindOne", user.Id).Return(user, nil)

	srv := mock.ServiceMock{}
	srv.On("UploadImage", file).Return("", err)

	hdr := NewHandler(&srv, &usrSrv, maxFileSize)

	hdr.UploadImage(&c)

	assert.Equal(t, want, c.V)
}

func (t *UserHandlerTest) TestUploadImageGrpcErr() {
	want := t.ServiceDownErr

	c := mock.ContextMock{}
	c.On("File", t.imgKey, constant.AllowImageContentType).Return(t.fileDecompose, nil)
	c.On("UserID").Return(t.user.Id)

	usrSrv := mockUsr.ServiceMock{}
	usrSrv.On("FindOne", t.user.Id).Return(nil, t.ServiceDownErr)

	srv := mock.ServiceMock{}
	srv.On("UploadImage", t.fileDecompose).Return("", nil)

	hdr := NewHandler(&srv, &usrSrv, t.maxFileSize)

	hdr.UploadImage(&c)

	assert.Equal(t.T(), want, c.V)
}
