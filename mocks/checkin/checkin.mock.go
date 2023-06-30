package checkin

import (
	"context"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/checkin/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) CheckinVerify(ctx context.Context, in *proto.CheckinVerifyRequest, opts ...grpc.CallOption) (res *proto.CheckinVerifyResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CheckinVerifyResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) CheckinConfirm(ctx context.Context, in *proto.CheckinConfirmRequest, opts ...grpc.CallOption) (res *proto.CheckinConfirmResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CheckinConfirmResponse)
	}

	return res, args.Error(1)
}

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) CheckinVerify(userid string, eventType int) (res *proto.CheckinVerifyResponse, err *dto.ResponseErr) {
	args := s.Called(userid, eventType)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CheckinVerifyResponse)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return res, err
}

func (s *ServiceMock) CheckinConfirm(token string) (res *proto.CheckinConfirmResponse, err *dto.ResponseErr) {
	args := s.Called(token)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CheckinConfirmResponse)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return res, err
}
