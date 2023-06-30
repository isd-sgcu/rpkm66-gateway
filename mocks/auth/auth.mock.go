package auth

import (
	"context"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/auth/auth/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) VerifyTicket(_ context.Context, in *proto.VerifyTicketRequest, _ ...grpc.CallOption) (res *proto.VerifyTicketResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.VerifyTicketResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Validate(_ context.Context, in *proto.ValidateRequest, _ ...grpc.CallOption) (res *proto.ValidateResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.ValidateResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) RefreshToken(_ context.Context, in *proto.RefreshTokenRequest, _ ...grpc.CallOption) (res *proto.RefreshTokenResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.RefreshTokenResponse)
	}

	return res, args.Error(1)
}

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) VerifyTicket(ticket string) (credential *proto.Credential, err *dto.ResponseErr) {
	args := s.Called(ticket)

	if args.Get(0) != nil {
		credential = args.Get(0).(*proto.Credential)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return credential, err
}

func (s *ServiceMock) Validate(token string) (payload *dto.TokenPayloadAuth, err *dto.ResponseErr) {
	args := s.Called(token)

	if args.Get(0) != nil {
		payload = args.Get(0).(*dto.TokenPayloadAuth)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return payload, err
}

func (s *ServiceMock) RefreshToken(token string) (credential *proto.Credential, err *dto.ResponseErr) {
	args := s.Called(token)

	if args.Get(0) != nil {
		credential = args.Get(0).(*proto.Credential)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return credential, err
}
