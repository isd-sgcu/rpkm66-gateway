package estamp

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

func (c *ClientMock) FindAllEvent(_ context.Context, in *proto.FindAllEventRequest, opts ...grpc.CallOption) (*proto.FindAllEventResponse, error) {
	args := c.Called(in)

	res := &proto.FindAllEventResponse{}

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindAllEventResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) FindEventByID(_ context.Context, in *proto.FindEventByIDRequest, opts ...grpc.CallOption) (*proto.FindEventByIDResponse, error) {
	args := c.Called(in)

	res := &proto.FindEventByIDResponse{}

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindEventByIDResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) FindAllEventWithType(_ context.Context, in *proto.FindAllEventWithTypeRequest, opts ...grpc.CallOption) (*proto.FindAllEventWithTypeResponse, error) {
	args := c.Called(in)

	res := &proto.FindAllEventWithTypeResponse{}

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindAllEventWithTypeResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Create(_ context.Context, in *proto.CreateEventRequest, opts ...grpc.CallOption) (*proto.CreateEventResponse, error) {
	return nil, nil
}

func (c *ClientMock) Update(_ context.Context, in *proto.UpdateEventRequest, opts ...grpc.CallOption) (*proto.UpdateEventResponse, error) {
	return nil, nil
}

func (c *ClientMock) Delete(_ context.Context, in *proto.DeleteEventRequest, opts ...grpc.CallOption) (*proto.DeleteEventResponse, error) {
	return nil, nil
}

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) FindEventByID(id string) (res *proto.FindEventByIDResponse, err *dto.ResponseErr) {
	args := s.Called(id)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindEventByIDResponse)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return res, err
}

func (s *ServiceMock) FindAllEventWithType(eventType string) (res *proto.FindAllEventWithTypeResponse, err *dto.ResponseErr) {
	args := s.Called(eventType)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindAllEventWithTypeResponse)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return res, err
}
