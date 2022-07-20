package group

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/validator"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/group"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type GroupHandlerTest struct {
	suite.Suite
	Group          *proto.Group
	GroupDto       *dto.GroupDto
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	InvalidIdErr   *dto.ResponseErr
	InvalidReqErr  *dto.ResponseErr
	ForbiddenErr   *dto.ResponseErr
	InternalErr    *dto.ResponseErr
}

func TestGroupHandler(t *testing.T) {
	suite.Run(t, new(GroupHandlerTest))
}

func (t *GroupHandlerTest) SetupTest() {
	t.Group = &proto.Group{
		Id:       faker.UUIDDigit(),
		LeaderID: faker.Word(),
		Token:    faker.Word(),
	}

	t.GroupDto = &dto.GroupDto{
		ID:       t.Group.Id,
		LeaderID: t.Group.LeaderID,
		Token:    t.Group.Token,
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
		Message:    "Invalid User ID",
		Data:       nil,
	}

	t.InvalidReqErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid Request Body",
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

func (t *GroupHandlerTest) TestFindOneSuccess() {
	want := t.Group

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID, nil)

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Group.LeaderID).Return(t.Group, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestFindOneNotFound() {
	want := t.NotFoundErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID, nil)

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Group.LeaderID).Return(nil, t.NotFoundErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID, nil)

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Group.LeaderID).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestFindByTokenSuccess() {
	want := &proto.FindByTokenGroupResponse{
		Id:    t.Group.Id,
		Token: t.Group.Token,
		Leader: &proto.UserInfo{
			Id:        faker.UUIDDigit(),
			FirstName: faker.Word(),
			LastName:  faker.Word(),
			ImageUrl:  faker.URL(),
		},
	}

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)

	srv := new(mock.ServiceMock)
	srv.On("FindByToken", t.Group.Token).Return(want, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindByToken(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestFindByTokenNotFound() {
	want := t.NotFoundErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)

	srv := new(mock.ServiceMock)
	srv.On("FindByToken", t.Group.Token).Return(nil, t.NotFoundErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindByToken(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestFindByTokenGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)

	srv := new(mock.ServiceMock)
	srv.On("FindByToken", t.Group.Token).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindByToken(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestUpdateSuccess() {
	want := t.Group

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.GroupDto{}).Return(t.GroupDto, nil)

	srv := new(mock.ServiceMock)
	srv.On("Update", t.GroupDto, t.Group.LeaderID).Return(t.Group, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestUpdateInvalidRequest() {
	want := t.InvalidReqErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.GroupDto{}).Return(nil, errors.New(t.InvalidReqErr.Message))

	srv := new(mock.ServiceMock)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestUpdateNotFound() {
	want := t.NotFoundErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.GroupDto{}).Return(t.GroupDto, nil)

	srv := new(mock.ServiceMock)
	srv.On("Update", t.GroupDto, t.Group.LeaderID).Return(nil, t.NotFoundErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestUpdateInvalidID() {
	want := t.InvalidIdErr

	c := &mock.ContextMock{}
	c.On("UserID").Return("abc")
	c.On("Bind", &dto.GroupDto{}).Return(t.GroupDto, nil)

	srv := new(mock.ServiceMock)
	srv.On("Update", t.GroupDto, "abc").Return(nil, t.InvalidIdErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestUpdateGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.GroupDto{}).Return(t.GroupDto, nil)

	srv := new(mock.ServiceMock)
	srv.On("Update", t.GroupDto, t.Group.LeaderID).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestJoinSuccess() {
	want := t.Group

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)
	c.On("UserID").Return(t.Group.LeaderID)
	srv := new(mock.ServiceMock)
	srv.On("Join", t.Group.Token, t.Group.LeaderID).Return(t.Group, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Join(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestJoinForbidden() {
	want := t.ForbiddenErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)
	c.On("UserID").Return(t.Group.LeaderID)
	srv := new(mock.ServiceMock)
	srv.On("Join", t.Group.Token, t.Group.LeaderID).Return(nil, t.ForbiddenErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Join(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestJoinInvalidId() {
	want := t.InvalidIdErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)
	c.On("UserID").Return("abc")
	srv := new(mock.ServiceMock)
	srv.On("Join", t.Group.Token, "abc").Return(nil, t.InvalidIdErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Join(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestJoinNotFound() {
	want := t.NotFoundErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)
	c.On("UserID").Return(t.Group.LeaderID)
	srv := new(mock.ServiceMock)
	srv.On("Join", t.Group.Token, t.Group.LeaderID).Return(nil, t.NotFoundErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Join(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestJoinGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)
	c.On("UserID").Return(t.Group.LeaderID)
	srv := new(mock.ServiceMock)
	srv.On("Join", t.Group.Token, t.Group.LeaderID).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Join(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestDeleteMemberSuccess() {
	want := t.Group

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.LeaderID, nil)
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("DeleteMember", t.Group.LeaderID, t.Group.LeaderID).Return(t.Group, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.DeleteMember(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestDeleteMemberNotFound() {
	want := t.NotFoundErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.LeaderID, nil)
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("DeleteMember", t.Group.LeaderID, t.Group.LeaderID).Return(nil, t.NotFoundErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.DeleteMember(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestDeleteMemberInvalidID() {
	want := t.InvalidIdErr

	c := &mock.ContextMock{}
	c.On("Param").Return("", errors.New(t.InvalidIdErr.Message))

	srv := new(mock.ServiceMock)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.DeleteMember(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestDeleteMemberForbidden() {
	want := t.ForbiddenErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.LeaderID, nil)
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("DeleteMember", t.Group.LeaderID, t.Group.LeaderID).Return(nil, t.ForbiddenErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.DeleteMember(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestDeleteMemberGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ContextMock{}
	c.On("Param").Return(t.Group.LeaderID, nil)
	c.On("UserID").Return(t.Group.LeaderID, nil)

	srv := new(mock.ServiceMock)
	srv.On("DeleteMember", t.Group.LeaderID, t.Group.LeaderID).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.DeleteMember(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestLeaveSuccess() {
	want := t.Group

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("Leave", t.Group.LeaderID).Return(t.Group, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Leave(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestLeaveInvalidId() {
	want := t.InvalidIdErr

	c := &mock.ContextMock{}
	c.On("UserID").Return("abc")

	srv := new(mock.ServiceMock)
	srv.On("Leave", "abc").Return(nil, t.InvalidIdErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Leave(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestLeaveNotFound() {
	want := t.NotFoundErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("Leave", t.Group.LeaderID).Return(nil, t.NotFoundErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Leave(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestLeaveForbidden() {
	want := t.ForbiddenErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("Leave", t.Group.LeaderID).Return(nil, t.ForbiddenErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Leave(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestLeaveInternalErr() {
	want := t.InternalErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("Leave", t.Group.LeaderID).Return(nil, t.InternalErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Leave(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestLeaveGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("Leave", t.Group.LeaderID).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Leave(c)

	assert.Equal(t.T(), want, c.V)
}
