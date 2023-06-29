package checkin

import (
	"context"

	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	"github.com/isd-sgcu/rpkm66-gateway/src/proto"
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

type ContextMock struct {
	mock.Mock
	V      interface{}
	Status int
}

func (c *ContextMock) JSON(status int, v interface{}) {
	c.V = v
	c.Status = status
}

func (c *ContextMock) Bind(v interface{}) error {
	args := c.Called(v)

	if args.Get(0) != nil {
		switch v.(type) {
		case *dto.UserDto:
			*v.(*dto.UserDto) = *args.Get(0).(*dto.UserDto)
		case *dto.Verify:
			*v.(*dto.Verify) = *args.Get(0).(*dto.Verify)
		case *dto.CheckinVerifyRequest:
			*v.(*dto.CheckinVerifyRequest) = *args.Get(0).(*dto.CheckinVerifyRequest)
		case *dto.CheckinConfirmRequest:
			*v.(*dto.CheckinConfirmRequest) = *args.Get(0).(*dto.CheckinConfirmRequest)
		}
	}

	return args.Error(1)
}

func (c *ContextMock) ID() (string, error) {
	args := c.Called()

	return args.String(0), args.Error(1)
}

func (c *ContextMock) Host() string {
	args := c.Called()

	return args.String(0)
}

func (c *ContextMock) UserID() string {
	args := c.Called()
	return args.String(0)
}
