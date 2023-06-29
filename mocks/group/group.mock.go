package group

import (
	"context"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) SelectBaan(userId string, baanIds []string) (result bool, err *dto.ResponseErr) {
	args := s.Called(userId, baanIds)

	if args.Get(0) != nil {
		result = args.Bool(0)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}
	return
}

func (s *ServiceMock) FindOne(id string) (result *proto.Group, err *dto.ResponseErr) {
	args := s.Called(id)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Group)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}
	return
}

func (s *ServiceMock) FindByToken(token string) (result *proto.FindByTokenGroupResponse, err *dto.ResponseErr) {
	args := s.Called(token)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.FindByTokenGroupResponse)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Create(id string) (result *proto.Group, err *dto.ResponseErr) {
	args := s.Called(id)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Group)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Update(in *dto.GroupDto, id string) (result *proto.Group, err *dto.ResponseErr) {
	args := s.Called(in, id)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Group)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Join(token string, userId string) (result *proto.Group, err *dto.ResponseErr) {
	args := s.Called(token, userId)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Group)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) DeleteMember(userId string, leaderId string) (result *proto.Group, err *dto.ResponseErr) {
	args := s.Called(userId, leaderId)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Group)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Leave(userId string) (result *proto.Group, err *dto.ResponseErr) {
	args := s.Called(userId)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Group)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) SelectBaan(_ context.Context, in *proto.SelectBaanRequest, _ ...grpc.CallOption) (result *proto.SelectBaanResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.SelectBaanResponse)
	}

	return result, args.Error(1)
}

func (c *ClientMock) FindOne(_ context.Context, in *proto.FindOneGroupRequest, _ ...grpc.CallOption) (res *proto.FindOneGroupResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindOneGroupResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) FindByToken(_ context.Context, in *proto.FindByTokenGroupRequest, _ ...grpc.CallOption) (res *proto.FindByTokenGroupResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindByTokenGroupResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Update(_ context.Context, in *proto.UpdateGroupRequest, _ ...grpc.CallOption) (res *proto.UpdateGroupResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.UpdateGroupResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Join(_ context.Context, in *proto.JoinGroupRequest, _ ...grpc.CallOption) (res *proto.JoinGroupResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.JoinGroupResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) DeleteMember(_ context.Context, in *proto.DeleteMemberGroupRequest, _ ...grpc.CallOption) (res *proto.DeleteMemberGroupResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.DeleteMemberGroupResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Leave(_ context.Context, in *proto.LeaveGroupRequest, _ ...grpc.CallOption) (res *proto.LeaveGroupResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.LeaveGroupResponse)
	}

	return res, args.Error(1)
}
