package file

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	"github.com/isd-sgcu/rpkm66-gateway/src/constant/file"
	mock "github.com/isd-sgcu/rpkm66-gateway/src/mocks/file"
	"github.com/isd-sgcu/rpkm66-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileServiceTest struct {
	suite.Suite
	url            string
	userId         string
	fileDecomposed *dto.DecomposedFile
	ServiceDownErr *dto.ResponseErr
}

func TestFileService(t *testing.T) {
	suite.Run(t, new(FileServiceTest))
}

func (t *FileServiceTest) SetupTest() {
	t.url = faker.URL()
	t.userId = faker.UUIDDigit()

	t.fileDecomposed = &dto.DecomposedFile{
		Filename: faker.Word(),
		Data:     []byte("Hello"),
	}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}
}

func (t *FileServiceTest) TestUploadSuccess() {
	want := t.url

	c := mock.ClientMock{}
	c.On("Upload", &proto.UploadRequest{
		Filename: t.fileDecomposed.Filename, Data: t.fileDecomposed.Data, Tag: 1, UserId: t.userId, Type: file.Image}).Return(&proto.UploadResponse{Url: t.url}, nil)

	srv := NewService(&c)

	actual, err := srv.Upload(t.fileDecomposed, t.userId, file.Profile, file.Image)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *FileServiceTest) TestUploadFailed() {
	want := t.ServiceDownErr

	c := mock.ClientMock{}
	c.On("Upload", &proto.UploadRequest{
		Filename: t.fileDecomposed.Filename, Data: t.fileDecomposed.Data, Tag: 1, UserId: t.userId, Type: file.Image}).Return(nil, errors.New("Cannot connect to service"))

	srv := NewService(&c)

	actual, err := srv.Upload(t.fileDecomposed, t.userId, file.Profile, file.Image)

	assert.Equal(t.T(), "", actual)
	assert.Equal(t.T(), want, err)
}
