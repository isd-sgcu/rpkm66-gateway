package baan

import (
	"context"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) FindAllBaan(_ context.Context, in *proto.FindAllBaanRequest, _ ...grpc.CallOption) (res *proto.FindAllBaanResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindAllBaanResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) FindOneBaan(_ context.Context, in *proto.FindOneBaanRequest, _ ...grpc.CallOption) (res *proto.FindOneBaanResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindOneBaanResponse)
	}

	return res, args.Error(1)
}

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) FindAll() (res []*proto.Baan, err *dto.ResponseErr) {
	args := s.Called()

	if args.Get(0) != nil {
		res = args.Get(0).([]*proto.Baan)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return res, err
}

func (s *ServiceMock) FindOne(id string) (res *proto.Baan, err *dto.ResponseErr) {
	args := s.Called(id)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.Baan)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return res, err
}
