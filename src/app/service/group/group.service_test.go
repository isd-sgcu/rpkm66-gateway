package group

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/utils"
	"github.com/isd-sgcu/rnkm65-gateway/src/mocks/group"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"testing"
)

type GroupServiceTest struct {
	suite.Suite
	User           *proto.User
	UserDto        *dto.UserDto
	Group          *proto.Group
	GroupDto       *dto.GroupDto
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	InvalidIdErr   *dto.ResponseErr
	ForbiddenErr   *dto.ResponseErr
	InternalErr    *dto.ResponseErr
}

func TestGroupService(t *testing.T) {
	suite.Run(t, new(GroupServiceTest))
}

func (t *GroupServiceTest) SetupTest() {
	t.User = &proto.User{
		Id:              faker.UUIDDigit(),
		Firstname:       faker.FirstName(),
		Lastname:        faker.LastName(),
		Nickname:        faker.Name(),
		Phone:           faker.Phonenumber(),
		LineID:          faker.Word(),
		Email:           faker.Email(),
		AllergyFood:     faker.Word(),
		FoodRestriction: faker.Word(),
		AllergyMedicine: faker.Word(),
		Disease:         faker.Word(),
		CanSelectBaan:   true,
		GroupId:         faker.UUIDDigit(),
	}

	t.UserDto = &dto.UserDto{
		ID:              t.User.Id,
		Firstname:       t.User.Firstname,
		Lastname:        t.User.Lastname,
		Nickname:        t.User.Nickname,
		Phone:           t.User.Phone,
		LineID:          t.User.LineID,
		Email:           t.User.Email,
		AllergyFood:     t.User.AllergyFood,
		FoodRestriction: t.User.FoodRestriction,
		AllergyMedicine: t.User.AllergyMedicine,
		Disease:         t.User.Disease,
		CanSelectBaan:   utils.BoolAdr(t.User.CanSelectBaan),
		GroupId:         t.User.GroupId,
	}

	t.Group = &proto.Group{
		Id:       t.User.GroupId,
		LeaderID: t.User.Id,
		Token:    faker.Word(),
		Members:  []*proto.User{t.User},
	}

	t.GroupDto = &dto.GroupDto{
		ID:       t.Group.Id,
		LeaderID: t.Group.LeaderID,
		Token:    t.Group.Token,
		Members:  []*dto.UserDto{t.UserDto},
	}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Group not found",
		Data:       nil,
	}

	t.InvalidIdErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid user id",
		Data:       nil,
	}

	t.ForbiddenErr = &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    "Not allowed",
		Data:       nil,
	}
	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Fail to create group",
		Data:       nil,
	}
}

func (t *GroupServiceTest) TestFindOneSuccess() {
	want := t.Group

	c := &group.ClientMock{}
	c.On("FindOne", &proto.FindOneGroupRequest{Id: t.User.Id}).Return(&proto.FindOneGroupResponse{Group: t.Group}, nil)

	srv := NewService(c)
	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestFindOneNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("FindOne", &proto.FindOneGroupRequest{Id: t.User.Id}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)
	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestFindOneInvalidId() {
	want := t.InvalidIdErr

	c := &group.ClientMock{}
	c.On("FindOne", &proto.FindOneGroupRequest{Id: "abc"}).Return(nil, status.Error(codes.InvalidArgument, "Invalid user id"))

	srv := NewService(c)
	actual, err := srv.FindOne("abc")

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("FindOne", &proto.FindOneGroupRequest{Id: t.User.Id}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)
	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestFindByTokenSuccess() {
	want := t.Group

	c := &group.ClientMock{}
	c.On("FindByToken", &proto.FindByTokenGroupRequest{Token: t.Group.Token}).Return(&proto.FindByTokenGroupResponse{Group: want}, nil)

	srv := NewService(c)
	actual, err := srv.FindByToken(t.Group.Token)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestFindByTokenNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("FindByToken", &proto.FindByTokenGroupRequest{Token: t.Group.Token}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)

	actual, err := srv.FindByToken(t.Group.Token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestFindByTokenGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("FindByToken", &proto.FindByTokenGroupRequest{Token: t.Group.Token}).Return(nil, errors.New("Server is down"))

	srv := NewService(c)

	actual, err := srv.FindByToken(t.Group.Token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestCreateSuccess() {
	want := t.Group

	c := &group.ClientMock{}
	c.On("Create", &proto.CreateGroupRequest{UserId: t.User.Id}).Return(&proto.CreateGroupResponse{Group: t.Group}, nil)

	srv := NewService(c)

	actual, err := srv.Create(t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestCreateNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("Create", &proto.CreateGroupRequest{UserId: t.User.Id}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)
	actual, err := srv.Create(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestCreateInvalidId() {
	want := t.InvalidIdErr

	c := &group.ClientMock{}
	c.On("Create", &proto.CreateGroupRequest{UserId: "abc"}).Return(nil, status.Error(codes.InvalidArgument, "Invalid user id"))

	srv := NewService(c)
	actual, err := srv.Create("abc")

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestCreateInternalErr() {
	want := t.InternalErr

	c := &group.ClientMock{}
	c.On("Create", &proto.CreateGroupRequest{UserId: t.User.Id}).Return(nil, status.Error(codes.Internal, "Fail to create group"))

	srv := NewService(c)
	actual, err := srv.Create(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestCreateGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("Create", &proto.CreateGroupRequest{UserId: t.User.Id}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.Create(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestUpdateSuccess() {
	want := t.Group

	c := &group.ClientMock{}
	c.On("Update", &proto.UpdateGroupRequest{Group: t.Group, LeaderId: t.Group.LeaderID}).Return(&proto.UpdateGroupResponse{Group: t.Group}, nil)

	srv := NewService(c)

	actual, err := srv.Update(t.GroupDto, t.Group.LeaderID)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestUpdateNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("Update", &proto.UpdateGroupRequest{Group: t.Group, LeaderId: t.Group.LeaderID}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)

	actual, err := srv.Update(t.GroupDto, t.Group.LeaderID)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestUpdateInvalidId() {
	want := t.InvalidIdErr

	c := &group.ClientMock{}
	c.On("Update", &proto.UpdateGroupRequest{Group: t.Group, LeaderId: "abc"}).Return(nil, status.Error(codes.InvalidArgument, "Invalid user id"))

	srv := NewService(c)

	actual, err := srv.Update(t.GroupDto, "abc")

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestUpdateGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("Update", &proto.UpdateGroupRequest{Group: t.Group, LeaderId: t.Group.LeaderID}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.Update(t.GroupDto, t.Group.LeaderID)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestJoinSuccess() {
	want := t.Group

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: t.User.Id, IsLeader: false, Members: 2}).Return(&proto.JoinGroupResponse{Group: t.Group}, nil)

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, t.User.Id, false, 2)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestJoinForbidden() {
	want := t.ForbiddenErr

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: t.User.Id, IsLeader: true, Members: 2}).Return(nil, status.Error(codes.PermissionDenied, "Not allowed"))

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, t.User.Id, true, 2)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestJoinInvalidId() {
	want := t.InvalidIdErr

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: "abc", IsLeader: false, Members: 2}).Return(nil, status.Error(codes.InvalidArgument, "Invalid user id"))

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, "abc", false, 2)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestJoinNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: t.User.Id, IsLeader: false, Members: 2}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, t.User.Id, false, 2)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestJoinGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: t.User.Id, IsLeader: false, Members: 2}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, t.User.Id, false, 2)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestDeleteMemberSuccess() {
	want := t.Group

	c := &group.ClientMock{}
	c.On("DeleteMember", &proto.DeleteMemberGroupRequest{UserId: t.User.Id, LeaderId: t.Group.LeaderID}).Return(&proto.DeleteMemberGroupResponse{Group: t.Group}, nil)

	srv := NewService(c)

	actual, err := srv.DeleteMember(t.User.Id, t.Group.LeaderID)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestDeleteMemberInvalidId() {
	want := t.InvalidIdErr

	c := &group.ClientMock{}
	c.On("DeleteMember", &proto.DeleteMemberGroupRequest{UserId: "abc", LeaderId: t.Group.LeaderID}).Return(nil, status.Error(codes.InvalidArgument, "Invalid user id"))

	srv := NewService(c)

	actual, err := srv.DeleteMember("abc", t.Group.LeaderID)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestDeleteMemberNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("DeleteMember", &proto.DeleteMemberGroupRequest{UserId: t.User.Id, LeaderId: t.Group.LeaderID}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)

	actual, err := srv.DeleteMember(t.User.Id, t.Group.LeaderID)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestDeleteMemberForbidden() {
	want := t.ForbiddenErr

	c := &group.ClientMock{}
	c.On("DeleteMember", &proto.DeleteMemberGroupRequest{UserId: t.User.Id, LeaderId: t.Group.LeaderID}).Return(nil, status.Error(codes.PermissionDenied, "Not allowed"))

	srv := NewService(c)

	actual, err := srv.DeleteMember(t.User.Id, t.Group.LeaderID)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestDeleteMemberGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("DeleteMember", &proto.DeleteMemberGroupRequest{UserId: t.User.Id, LeaderId: t.Group.LeaderID}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.DeleteMember(t.User.Id, t.Group.LeaderID)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestLeaveSuccess() {
	want := t.Group

	c := &group.ClientMock{}
	c.On("Leave", &proto.LeaveGroupRequest{UserId: t.User.Id}).Return(&proto.LeaveGroupResponse{Group: t.Group}, nil)

	srv := NewService(c)

	actual, err := srv.Leave(t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestLeaveInvalidId() {
	want := t.InvalidIdErr

	c := &group.ClientMock{}
	c.On("Leave", &proto.LeaveGroupRequest{UserId: "abc"}).Return(nil, status.Error(codes.InvalidArgument, "Invalid user id"))

	srv := NewService(c)

	actual, err := srv.Leave("abc")

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestLeaveNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("Leave", &proto.LeaveGroupRequest{UserId: t.User.Id}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)

	actual, err := srv.Leave(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestLeaveForbidden() {
	want := t.ForbiddenErr

	c := &group.ClientMock{}
	c.On("Leave", &proto.LeaveGroupRequest{UserId: t.User.Id}).Return(nil, status.Error(codes.PermissionDenied, "Not allowed"))

	srv := NewService(c)

	actual, err := srv.Leave(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestLeaveInternalErr() {
	want := t.InternalErr

	c := &group.ClientMock{}
	c.On("Leave", &proto.LeaveGroupRequest{UserId: t.User.Id}).Return(nil, status.Error(codes.Internal, "Fail to create group"))

	srv := NewService(c)

	actual, err := srv.Leave(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestLeaveGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("Leave", &proto.LeaveGroupRequest{UserId: t.User.Id}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.Leave(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}
