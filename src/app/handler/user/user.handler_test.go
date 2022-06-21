package user

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/validator"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/user"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type UserHandlerTest struct {
	suite.Suite
	User           *proto.User
	UserDto        *dto.UserDto
	BindErr        *dto.ResponseErr
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(UserHandlerTest))
}

func (t *UserHandlerTest) SetupTest() {
	t.User = &proto.User{
		Id:                    faker.UUIDDigit(),
		Firstname:             faker.FirstName(),
		Lastname:              faker.LastName(),
		Nickname:              faker.Name(),
		StudentID:             faker.Word(),
		Faculty:               faker.Word(),
		Year:                  faker.Word(),
		Phone:                 faker.Phonenumber(),
		LineID:                faker.Word(),
		Email:                 faker.Email(),
		AllergyFood:           faker.Word(),
		FoodRestriction:       faker.Word(),
		AllergyMedicine:       faker.Word(),
		Disease:               faker.Word(),
		VaccineCertificateUrl: faker.URL(),
		ImageUrl:              faker.URL(),
	}

	t.UserDto = &dto.UserDto{
		ID:                    t.User.Id,
		Firstname:             t.User.Firstname,
		Lastname:              t.User.Lastname,
		Nickname:              t.User.Nickname,
		StudentID:             t.User.StudentID,
		Faculty:               t.User.Faculty,
		Year:                  t.User.Year,
		Phone:                 t.User.Phone,
		LineID:                t.User.LineID,
		Email:                 t.User.Email,
		AllergyFood:           t.User.AllergyFood,
		FoodRestriction:       t.User.FoodRestriction,
		AllergyMedicine:       t.User.AllergyMedicine,
		Disease:               t.User.Disease,
		VaccineCertificateUrl: t.User.VaccineCertificateUrl,
		ImageUrl:              t.User.ImageUrl,
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
}

func (t *UserHandlerTest) TestFindOneUser() {
	want := t.User

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.User.Id).Return(want, nil)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("UserID").Return(t.User.Id)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestFindOneFoundErr() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.User.Id).Return(nil, t.NotFoundErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("UserID").Return(t.User.Id)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.User.Id).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("UserID").Return(t.User.Id)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestCreateSuccess() {
	want := t.User

	srv := new(mock.ServiceMock)
	srv.On("Create", t.UserDto).Return(want, nil)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("Bind", &dto.UserDto{}).Return(nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Create(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestCreateValidateErr() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid body request",
		Data: []*dto.BadReqErrResponse{
			&dto.BadReqErrResponse{
				Message:     "VaccineCertificateUrl must be a valid URL",
				FailedField: "VaccineCertificateUrl",
				Value:       "",
			},
		},
	}

	srv := new(mock.ServiceMock)
	srv.On("Create", t.UserDto).Return(t.User, nil)

	t.UserDto.VaccineCertificateUrl = ""

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("Bind", &dto.UserDto{}).Return(nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Create(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestCreateGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("Create", t.UserDto).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("Bind", &dto.UserDto{}).Return(nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Create(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUpdateSuccess() {
	want := t.User

	srv := new(mock.ServiceMock)
	srv.On("Update", t.User.Id, t.UserDto).Return(want, nil)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("ID").Return(t.User.Id, nil)
	c.On("Bind", &dto.UserDto{}).Return(nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUpdateValidateErr() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid body request",
		Data: []*dto.BadReqErrResponse{
			&dto.BadReqErrResponse{
				Message:     "VaccineCertificateUrl must be a valid URL",
				FailedField: "VaccineCertificateUrl",
				Value:       "",
			},
		},
	}

	t.UserDto.VaccineCertificateUrl = ""

	srv := new(mock.ServiceMock)
	srv.On("Update", t.UserDto).Return(nil, t.BindErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("Bind", &dto.UserDto{}).Return(nil)
	c.On("ID").Return(t.User.Id, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUpdateNotFound() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("Update", t.User.Id, t.UserDto).Return(nil, t.NotFoundErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("ID").Return(t.User.Id, nil)
	c.On("Bind", &dto.UserDto{}).Return(nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUpdateInvalidID() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "ID must be the uuid",
	}

	srv := new(mock.ServiceMock)
	srv.On("Update", t.User.Id).Return(nil, t.NotFoundErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("ID").Return("", errors.New(want.Message))

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestUpdateGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("Update", t.User.Id, t.UserDto).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("ID").Return(t.User.Id, nil)
	c.On("Bind", &dto.UserDto{}).Return(nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Update(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestCreateOrUpdateSuccess() {
	want := t.User

	srv := new(mock.ServiceMock)
	srv.On("CreateOrUpdate", t.UserDto).Return(want, nil)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.UserDto{}).Return(nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.CreateOrUpdate(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestCreateOrUpdateValidateErr() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid body request",
		Data: []*dto.BadReqErrResponse{
			&dto.BadReqErrResponse{
				Message:     "ID is not uuid",
				FailedField: "ID",
				Value:       "abc",
			},
		},
	}

	t.User.Id = "abc"

	srv := new(mock.ServiceMock)
	srv.On("CreateOrUpdate", t.UserDto).Return(nil, t.BindErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("Bind", &dto.UserDto{}).Return(nil)
	c.On("UserID").Return(t.User.Id)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.CreateOrUpdate(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestCreateOrUpdateGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("CreateOrUpdate", t.UserDto).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.UserDto{}).Return(nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.CreateOrUpdate(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestDeleteSuccess() {
	srv := new(mock.ServiceMock)
	srv.On("Delete", t.User.Id).Return(true, nil)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("ID").Return(t.User.Id, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Delete(c)

	assert.True(t.T(), c.V.(bool))
}

func (t *UserHandlerTest) TestDeleteNotFound() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("Delete", t.User.Id).Return(false, t.NotFoundErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
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

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("ID").Return("", errors.New(want.Message))

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Delete(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *UserHandlerTest) TestDeleteGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("Delete", t.User.Id).Return(false, t.ServiceDownErr)

	c := &mock.ContextMock{
		User:    t.User,
		UserDto: t.UserDto,
	}
	c.On("ID").Return(t.User.Id, nil)
	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Delete(c)

	assert.Equal(t.T(), want, c.V)
}
