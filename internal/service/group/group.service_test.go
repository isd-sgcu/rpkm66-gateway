package group

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	constant "github.com/isd-sgcu/rpkm66-gateway/constant/baan"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/group"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GroupServiceTest struct {
	suite.Suite
	User                     *proto.UserInfo
	UserDto                  *dto.UserInfo
	Group                    *proto.Group
	GroupDto                 *dto.GroupDto
	FindByTokenGroupResponse *proto.FindByTokenGroupResponse
	NotFoundErr              *dto.ResponseErr
	ServiceDownErr           *dto.ResponseErr
	InvalidIdErr             *dto.ResponseErr
	ForbiddenErr             *dto.ResponseErr
	InternalErr              *dto.ResponseErr
}

func TestGroupService(t *testing.T) {
	suite.Run(t, new(GroupServiceTest))
}

func (t *GroupServiceTest) SetupTest() {
	t.User = &proto.UserInfo{
		Id:        faker.UUIDDigit(),
		Firstname: faker.FirstName(),
		Lastname:  faker.LastName(),
		ImageUrl:  faker.URL(),
	}

	t.UserDto = &dto.UserInfo{
		ID:        t.User.Id,
		Firstname: t.User.Firstname,
		Lastname:  t.User.Lastname,
		ImageUrl:  t.User.ImageUrl,
	}

	t.Group = &proto.Group{
		Id:       faker.UUIDDigit(),
		LeaderID: t.User.Id,
		Token:    faker.Word(),
		Members:  []*proto.UserInfo{t.User},
	}

	t.GroupDto = &dto.GroupDto{
		ID:       t.Group.Id,
		LeaderID: t.Group.LeaderID,
		Token:    t.Group.Token,
		Members:  []*dto.UserInfo{t.UserDto},
	}

	t.FindByTokenGroupResponse = &proto.FindByTokenGroupResponse{
		Id:     t.Group.Id,
		Token:  t.Group.Token,
		Leader: t.User,
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
	c.On("FindOne", &proto.FindOneGroupRequest{UserId: t.User.Id}).Return(&proto.FindOneGroupResponse{Group: t.Group}, nil)

	srv := NewService(c)
	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestFindOneNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("FindOne", &proto.FindOneGroupRequest{UserId: t.User.Id}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)
	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestFindOneInvalidId() {
	want := t.InvalidIdErr

	c := &group.ClientMock{}
	c.On("FindOne", &proto.FindOneGroupRequest{UserId: "abc"}).Return(nil, status.Error(codes.InvalidArgument, "Invalid user id"))

	srv := NewService(c)
	actual, err := srv.FindOne("abc")

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("FindOne", &proto.FindOneGroupRequest{UserId: t.User.Id}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)
	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestFindByTokenSuccess() {
	want := t.FindByTokenGroupResponse

	c := &group.ClientMock{}
	c.On("FindByToken", &proto.FindByTokenGroupRequest{Token: t.Group.Token}).Return(&proto.FindByTokenGroupResponse{Id: t.Group.Id, Token: t.Group.Token, Leader: t.User}, nil)

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
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: t.User.Id}).Return(&proto.JoinGroupResponse{Group: t.Group}, nil)

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *GroupServiceTest) TestJoinForbidden() {
	want := t.ForbiddenErr

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: t.User.Id}).Return(nil, status.Error(codes.PermissionDenied, "Not allowed"))

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestJoinInvalidId() {
	want := t.InvalidIdErr

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: "abc"}).Return(nil, status.Error(codes.InvalidArgument, "Invalid user id"))

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, "abc")

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestJoinNotFound() {
	want := t.NotFoundErr

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: t.User.Id}).Return(nil, status.Error(codes.NotFound, "Group not found"))

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestJoinGrpcErr() {
	want := t.ServiceDownErr

	c := &group.ClientMock{}
	c.On("Join", &proto.JoinGroupRequest{Token: t.Group.Token, UserId: t.User.Id}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.Join(t.Group.Token, t.User.Id)

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

func createBaanSlices() ([]*proto.Baan, []string) {
	var result []*proto.Baan
	var baanIds []string

	for i := 0; i < 3; i++ {
		b := &proto.Baan{
			Id:            uuid.New().String(),
			NameTH:        faker.Word(),
			DescriptionTH: faker.Paragraph(),
			NameEN:        faker.Word(),
			DescriptionEN: faker.Paragraph(),
			Size:          constant.M,
			Facebook:      faker.Word(),
			FacebookUrl:   faker.URL(),
			Instagram:     faker.Word(),
			InstagramUrl:  faker.URL(),
			Line:          faker.Word(),
			LineUrl:       faker.URL(),
			ImageUrl:      faker.URL(),
		}

		result = append(result, b)
		baanIds = append(baanIds, b.Id)
	}
	return result, baanIds
}

func (t *GroupServiceTest) TestUpdateSelectBaanSuccess() {
	_, baanIds := createBaanSlices()

	c := &group.ClientMock{}
	c.On("SelectBaan", &proto.SelectBaanRequest{
		UserId: t.Group.Id,
		Baans:  baanIds,
	}).Return(&proto.SelectBaanResponse{Success: true}, nil)

	srv := NewService(c)

	actual, err := srv.SelectBaan(t.Group.Id, baanIds)

	assert.Nil(t.T(), err)
	assert.True(t.T(), actual)
}

func (t *GroupServiceTest) TestUpdateSelectBaanInvalidInput() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid numbers of baan or bann is duplicated",
		Data:       nil,
	}

	_, baanIds := createBaanSlices()

	c := &group.ClientMock{}
	c.On("SelectBaan", &proto.SelectBaanRequest{
		UserId: t.Group.Id,
		Baans:  baanIds,
	}).Return(nil, status.Error(codes.InvalidArgument, "Duplicated baan"))

	srv := NewService(c)

	actual, err := srv.SelectBaan(t.Group.Id, baanIds)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestUpdateSelectBaanForbbidenInput() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    "Forbidden Action",
		Data:       nil,
	}

	_, baanIds := createBaanSlices()

	c := &group.ClientMock{}
	c.On("SelectBaan", &proto.SelectBaanRequest{
		UserId: t.Group.Id,
		Baans:  baanIds,
	}).Return(nil, status.Error(codes.PermissionDenied, "forbidden action"))

	srv := NewService(c)

	actual, err := srv.SelectBaan(t.Group.Id, baanIds)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *GroupServiceTest) TestUpdateSelectBaanInternalServiceErr() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Service is down",
		Data:       nil,
	}

	_, baanIds := createBaanSlices()

	c := &group.ClientMock{}
	c.On("SelectBaan", &proto.SelectBaanRequest{
		UserId: t.Group.Id,
		Baans:  baanIds,
	}).Return(nil, status.Error(codes.Internal, "Something wrong"))

	srv := NewService(c)

	actual, err := srv.SelectBaan(t.Group.Id, baanIds)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), want, err)
}
