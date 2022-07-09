package file

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/file"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type FileServiceTest struct {
	suite.Suite
	fileDecomposed *dto.DecomposedFile
	ServiceDownErr *dto.ResponseErr
}

func TestFileService(t *testing.T) {
	suite.Run(t, new(FileServiceTest))
}

func (t *FileServiceTest) SetupTest() {
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

func (t *FileServiceTest) TestUploadImageSuccess() {
	want := t.fileDecomposed.Filename

	c := mock.ClientMock{}
	c.On("UploadImage", &proto.UploadImageRequest{
		Filename: t.fileDecomposed.Filename, Data: t.fileDecomposed.Data}).Return(&proto.UploadImageResponse{Filename: t.fileDecomposed.Filename}, nil)

	srv := NewService(&c)

	actual, err := srv.UploadImage(t.fileDecomposed)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *FileServiceTest) TestUploadImageFailed() {
	want := t.ServiceDownErr

	c := mock.ClientMock{}
	c.On("UploadImage", &proto.UploadImageRequest{
		Filename: t.fileDecomposed.Filename, Data: t.fileDecomposed.Data}).Return(nil, errors.New("Cannot connect to service"))

	srv := NewService(&c)

	actual, err := srv.UploadImage(t.fileDecomposed)

	assert.Equal(t.T(), "", actual)
	assert.Equal(t.T(), want, err)
}
