package file

import (
	"context"

	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	"github.com/isd-sgcu/rpkm66-gateway/src/constant/file"
	"github.com/isd-sgcu/rpkm66-gateway/src/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) Upload(_ context.Context, in *proto.UploadRequest, _ ...grpc.CallOption) (res *proto.UploadResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.UploadResponse)
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

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) Upload(file *dto.DecomposedFile, userId string, tag file.Tag, fileType file.Type) (res string, err *dto.ResponseErr) {
	args := s.Called(file, userId, tag, fileType)

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return args.String(0), err
}
