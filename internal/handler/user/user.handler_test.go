package user

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/utils"
	"github.com/isd-sgcu/rpkm66-gateway/internal/validator"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/rctx"
	mock "github.com/isd-sgcu/rpkm66-gateway/mocks/user"
	eventProto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/event/v1"
	userProto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/user/v1"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTest struct {
	suite.Suite
	User           *userProto.User
	Events         []*eventProto.Event
	UserDto        *dto.UserDto
	UpdateUserDto  *dto.UpdateUserDto
	BindErr        *dto.ResponseErr
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	InternalErr    *dto.ResponseErr
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(UserHandlerTest))
}

func (t *UserHandlerTest) SetupTest() {
	t.User = &userProto.User{
		Id:              faker.UUIDDigit(),
		Title:           faker.Word(),
		Firstname:       faker.FirstName(),
		Lastname:        faker.LastName(),
		Nickname:        faker.Name(),
		StudentID:       faker.Word(),
		Faculty:         faker.Word(),
		Year:            faker.Word(),
		Phone:           faker.Phonenumber(),
		LineID:          faker.Word(),
		Email:           faker.Email(),
		AllergyFood:     faker.Word(),
		FoodRestriction: faker.Word(),
		AllergyMedicine: faker.Word(),
		Disease:         faker.Word(),
		ImageUrl:        faker.URL(),
		CanSelectBaan:   true,
		BaanId:          faker.UUIDDigit(),
	}

	t.Events = make([]*eventProto.Event, 3)

	t.Events[0] = &eventProto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.Events[1] = &eventProto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.Events[2] = &eventProto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.UserDto = &dto.UserDto{
		ID:              t.User.Id,
		Title:           t.User.Title,
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
	}

	t.UpdateUserDto = &dto.UpdateUserDto{
		Title:           t.User.Title,
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
	}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "User not found",
		Data:       nil,
	}

	t.BindErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid ID",
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		Data:       nil,
	}
}

func (t *UserHandlerTest) TestFindOneUser() {
	want := t.User

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.User.Id).Return(want, nil)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusOK, c.Status)
}

func (t *UserHandlerTest) TestFindOneFoundErr() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.User.Id).Return(nil, t.NotFoundErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusNotFound, c.Status)
}

func (t *UserHandlerTest) TestFindOneInternalErr() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Invalid ID",
		Data:       nil,
	}

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.User.Id).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return("", errors.New("Cannot parse id"))

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusInternalServerError, c.Status)
}

func (t *UserHandlerTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.User.Id).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

func (t *UserHandlerTest) TestCreateSuccess() {
	want := t.User

	srv := new(mock.ServiceMock)
	srv.On("Create", t.UserDto).Return(want, nil)

	c := &rctx.ContextMock{}
	c.On("Bind", &dto.UserDto{}).Return(t.UserDto, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Create(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusCreated, c.Status)
}

func (t *UserHandlerTest) TestCreateValidateErr() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid body request",
		Data: []*dto.BadReqErrResponse{
			{
				Message:     "Email must be a valid email address",
				FailedField: "Email",
				Value:       "",
			},
		},
	}

	t.UserDto.Email = ""

	srv := new(mock.ServiceMock)
	srv.On("Create", t.UserDto).Return(t.User, nil)

	c := &rctx.ContextMock{}
	c.On("Bind", &dto.UserDto{}).Return(t.UserDto, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Create(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *UserHandlerTest) TestCreateGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("Create", t.UserDto).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("Bind", &dto.UserDto{}).Return(t.UserDto, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Create(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

func (t *UserHandlerTest) TestUpdateSuccess() {
	want := t.User

	srv := new(mock.ServiceMock)
	srv.On("Update", t.User.Id, t.UpdateUserDto).Return(want, nil)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)
	c.On("UserID").Return(t.User.Id, nil)
	c.On("Bind", &dto.UpdateUserDto{}).Return(t.UpdateUserDto, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUpdateNotFound() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("Update", t.User.Id, t.UpdateUserDto).Return(nil, t.NotFoundErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)
	c.On("UserID").Return(t.User.Id, nil)
	c.On("Bind", &dto.UpdateUserDto{}).Return(t.UpdateUserDto, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusNotFound, c.Status)
}

func (t *UserHandlerTest) TestUpdateGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("Update", t.User.Id, t.UpdateUserDto).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)
	c.On("UserID").Return(t.User.Id, nil)
	c.On("Bind", &dto.UpdateUserDto{}).Return(t.UpdateUserDto, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

func (t *UserHandlerTest) TestCreateOrUpdateSuccess() {
	want := t.User

	srv := new(mock.ServiceMock)
	srv.On("CreateOrUpdate", t.UserDto).Return(want, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.UserDto{}).Return(t.UserDto, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.CreateOrUpdate(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusOK, c.Status)
}

func (t *UserHandlerTest) TestCreateOrUpdateValidateErr() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid body request",
		Data: []*dto.BadReqErrResponse{
			{
				Message:     "ID is not uuid",
				FailedField: "ID",
				Value:       "abc",
			},
		},
	}

	t.User.Id = "abc"

	srv := new(mock.ServiceMock)
	srv.On("CreateOrUpdate", t.UserDto).Return(nil, t.BindErr)

	c := &rctx.ContextMock{}
	c.On("Bind", &dto.UserDto{}).Return(t.UserDto, nil)
	c.On("UserID").Return(t.User.Id)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.CreateOrUpdate(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *UserHandlerTest) TestCreateOrUpdateGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("CreateOrUpdate", t.UserDto).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.UserDto{}).Return(t.UserDto, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.CreateOrUpdate(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

func (t *UserHandlerTest) TestDeleteSuccess() {
	srv := new(mock.ServiceMock)
	srv.On("Delete", t.User.Id).Return(true, nil)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Delete(c)

	assert.True(t.T(), c.V.(bool))
	assert.Equal(t.T(), http.StatusOK, c.Status)
}

func (t *UserHandlerTest) TestDeleteNotFound() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("Delete", t.User.Id).Return(false, t.NotFoundErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Delete(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestDeleteInvalidID() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "ID must be the uuid",
	}

	srv := new(mock.ServiceMock)
	srv.On("Delete", t.User.Id).Return(false, t.NotFoundErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return("", errors.New(want.Message))

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Delete(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *UserHandlerTest) TestDeleteGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("Delete", t.User.Id).Return(false, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.User.Id, nil)
	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Delete(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

func (t *UserHandlerTest) TestGetUserEstampSuccess() {
	want := &userProto.GetUserEstampResponse{
		EventList: []*eventProto.Event{
			t.Events[0],
			t.Events[1],
		},
	}

	s := &mock.ServiceMock{}
	s.On("GetUserEstamp", t.User.Id).Return(want, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id, nil)

	hdr.GetUserEstamp(c)

	assert.Equal(t.T(), http.StatusOK, c.Status)
	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestFindUserInternal() {
	s := &mock.ServiceMock{}
	s.On("GetUserEstamp", t.User.Id).Return(nil, t.InternalErr)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id, nil)

	hdr.GetUserEstamp(c)

	assert.Equal(t.T(), http.StatusInternalServerError, c.Status)
	assert.Equal(t.T(), t.InternalErr, c.V)
}

func (t *UserHandlerTest) TestFindUserUnavailable() {
	s := &mock.ServiceMock{}
	s.On("GetUserEstamp", t.User.Id).Return(nil, t.ServiceDownErr)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id, nil)

	hdr.GetUserEstamp(c)

	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
	assert.Equal(t.T(), t.ServiceDownErr, c.V)
}

func (t *UserHandlerTest) TestConfirmEstampSuccess() {
	want := &userProto.ConfirmEstampResponse{}

	s := &mock.ServiceMock{}
	s.On("ConfirmEstamp", t.User.Id, t.Events[0].Id).Return(want, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.ConfirmEstampRequest{}).Return(&dto.ConfirmEstampRequest{
		EventId: t.Events[0].Id,
	}, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.ConfirmEstamp(c)

	assert.Equal(t.T(), http.StatusNoContent, c.Status)
	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestConfirmEstampBadRequest() {
	s := &mock.ServiceMock{}

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.ConfirmEstampRequest{}).Return(nil, errors.New(""))

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.ConfirmEstamp(c)

	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *UserHandlerTest) TestConfirmEstampInnerError() {
	s := &mock.ServiceMock{}
	s.On("ConfirmEstamp", t.User.Id, t.Events[0].Id).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.ConfirmEstampRequest{}).Return(&dto.ConfirmEstampRequest{
		EventId: t.Events[0].Id,
	}, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.ConfirmEstamp(c)

	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}
