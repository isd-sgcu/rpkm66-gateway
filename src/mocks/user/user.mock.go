package user

import (
	"context"
	"github.com/isd-sgcu/rnkm65-gateway/src/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) FindOne(id string) (result *proto.User, err *dto.ResponseErr) {
	args := s.Called(id)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.User)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Create(in *dto.UserDto) (result *proto.User, err *dto.ResponseErr) {
	args := s.Called(in)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.User)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Update(id string, in *dto.UserDto) (result *proto.User, err *dto.ResponseErr) {
	args := s.Called(id, in)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.User)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) CreateOrUpdate(in *dto.UserDto) (result *proto.User, err *dto.ResponseErr) {
	args := s.Called(in)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.User)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Delete(id string) (result bool, err *dto.ResponseErr) {
	args := s.Called(id)

	result = args.Bool(0)

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) FindOne(_ context.Context, in *proto.FindOneUserRequest, _ ...grpc.CallOption) (res *proto.FindOneUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindOneUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Create(_ context.Context, in *proto.CreateUserRequest, _ ...grpc.CallOption) (res *proto.CreateUserResponse, err error) {
	args := c.Called(in.User)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CreateUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Update(_ context.Context, in *proto.UpdateUserRequest, _ ...grpc.CallOption) (res *proto.UpdateUserResponse, err error) {
	args := c.Called(in.User)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.UpdateUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Delete(_ context.Context, in *proto.DeleteUserRequest, _ ...grpc.CallOption) (res *proto.DeleteUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.DeleteUserResponse)
	}

	return res, args.Error(1)
}

func (s *ClientMock) CreateOrUpdate(_ context.Context, in *proto.CreateOrUpdateUserRequest, _ ...grpc.CallOption) (res *proto.CreateOrUpdateUserResponse, err error) {
	args := s.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CreateOrUpdateUserResponse)
	}

	return res, args.Error(1)
}

type ContextMock struct {
	mock.Mock
	V       interface{}
	User    *proto.User
	UserDto *dto.UserDto
}

func (c *ContextMock) JSON(_ int, v interface{}) {
	c.V = v
}

func (c *ContextMock) Bind(v interface{}) error {
	args := c.Called(v)

	*v.(*dto.UserDto) = *c.UserDto

	return args.Error(0)
}

func (c *ContextMock) ID() (string, error) {
	args := c.Called()

	return args.String(0), args.Error(1)
}

func (c *ContextMock) UserID() string {
	args := c.Called()
	return args.String(0)
}
