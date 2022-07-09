package file

import (
	"context"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) UploadImage(_ context.Context, in *proto.UploadImageRequest, _ ...grpc.CallOption) (res *proto.UploadImageResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.UploadImageResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) UploadFile(_ context.Context, in *proto.UploadFileRequest, _ ...grpc.CallOption) (res *proto.UploadFileResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.UploadFileResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) GetSignedUrl(_ context.Context, in *proto.GetSignedUrlRequest, _ ...grpc.CallOption) (res *proto.GetSignedUrlResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.GetSignedUrlResponse)
	}

	return res, args.Error(1)
}

type ContextMock struct {
	mock.Mock
	V interface{}
}

func (c *ContextMock) UserID() string {
	args := c.Called()
	return args.String(0)
}

func (c *ContextMock) JSON(_ int, v interface{}) {
	c.V = v
}

func (c *ContextMock) File(key string, allowContent map[string]struct{}, _ int64) (res *dto.DecomposedFile, err error) {
	args := c.Called(key, allowContent)

	if args.Get(0) != nil {
		res = args.Get(0).(*dto.DecomposedFile)
	}

	return res, args.Error(1)
}

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) UploadImage(file *dto.DecomposedFile) (res string, err *dto.ResponseErr) {
	args := s.Called(file)

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return args.String(0), err
}
