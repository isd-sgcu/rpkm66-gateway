package group

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/validator"
	mock "github.com/isd-sgcu/rpkm66-gateway/mocks/group"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/rctx"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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
			Firstname: faker.Word(),
			Lastname:  faker.Word(),
			ImageUrl:  faker.URL(),
		},
	}

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
	c.On("Param").Return(t.Group.Token, nil)

	srv := new(mock.ServiceMock)
	srv.On("FindByToken", t.Group.Token).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindByToken(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestJoinSuccess() {
	want := t.Group

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
	c.On("Param").Return(t.Group.LeaderID, nil)
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("DeleteMember", t.Group.LeaderID, t.Group.LeaderID).Return(nil, t.NotFoundErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.DeleteMember(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestDeleteMemberForbidden() {
	want := t.ForbiddenErr

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
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

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)

	srv := new(mock.ServiceMock)
	srv.On("Leave", t.Group.LeaderID).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Leave(c)

	assert.Equal(t.T(), want, c.V)
}

func createBaanSlices() *dto.SelectBaan {
	var baanIds []string
	for i := 0; i < 3; i++ {
		baanIds = append(baanIds, uuid.New().String())
	}

	return &dto.SelectBaan{Baans: baanIds}
}

func (t *GroupHandlerTest) TestSelectBaanSuccess() {
	baans := createBaanSlices()

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.SelectBaan{}).Return(baans, nil)

	srv := new(mock.ServiceMock)
	srv.On("SelectBaan", t.Group.LeaderID, baans.Baans).Return(true, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.SelectBaan(c)

	assert.Equal(t.T(), http.StatusNoContent, c.Status)
}

func (t *GroupHandlerTest) TestSelectBaanInvalidInput() {
	baans := createBaanSlices()

	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid Input",
		Data:       nil,
	}

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.SelectBaan{}).Return(baans, nil)

	srv := new(mock.ServiceMock)
	srv.On("SelectBaan", t.Group.LeaderID, baans.Baans).Return(false, want)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.SelectBaan(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestSelectBaanForbiddenActio() {
	baans := createBaanSlices()

	want := &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    "Forbidden Action",
		Data:       nil,
	}

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.SelectBaan{}).Return(baans, nil)

	srv := new(mock.ServiceMock)
	srv.On("SelectBaan", t.Group.LeaderID, baans.Baans).Return(false, want)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.SelectBaan(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestSelectBaanInternalErr() {
	baans := createBaanSlices()

	want := t.ServiceDownErr

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.SelectBaan{}).Return(baans, nil)

	srv := new(mock.ServiceMock)
	srv.On("SelectBaan", t.Group.LeaderID, baans.Baans).Return(false, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.SelectBaan(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *GroupHandlerTest) TestSelectBaanFailed() {
	baans := createBaanSlices()

	want := t.ServiceDownErr

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.Group.LeaderID)
	c.On("Bind", &dto.SelectBaan{}).Return(baans, nil)

	srv := new(mock.ServiceMock)
	srv.On("SelectBaan", t.Group.LeaderID, baans.Baans).Return(false, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.SelectBaan(c)

	assert.Equal(t.T(), want, c.V)
}
