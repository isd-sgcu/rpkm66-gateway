package auth

import (
	"context"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ContextMock struct {
	mock.Mock
	V               interface{}
	VerifyTicketDto *dto.VerifyTicket
	RefreshTokenDto *dto.RedeemNewToken
}

func (c *ContextMock) Bind(v interface{}) error {
	args := c.Called(v)

	switch v.(type) {
	case *dto.VerifyTicket:
		*v.(*dto.VerifyTicket) = *c.VerifyTicketDto
	case *dto.RedeemNewToken:
		*v.(*dto.RedeemNewToken) = *c.RefreshTokenDto
	}

	return args.Error(1)
}

func (c *ContextMock) JSON(_ int, v interface{}) {
	c.V = v
}

func (c *ContextMock) UserID() string {
	args := c.Called()

	return args.String(0)
}

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
