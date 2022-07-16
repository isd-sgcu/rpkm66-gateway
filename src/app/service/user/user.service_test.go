package user

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/utils"
	"github.com/isd-sgcu/rnkm65-gateway/src/mocks/user"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"testing"
)

type UserServiceTest struct {
	suite.Suite
	User           *proto.User
	UserReq        *proto.User
	UserDto        *dto.UserDto
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTest))
}

func (t *UserServiceTest) SetupTest() {
	t.User = &proto.User{
		Id:              faker.UUIDDigit(),
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
		CanSelectBaan:   true,
	}

	t.UserReq = &proto.User{
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
		ImageUrl:        t.User.ImageUrl,
		CanSelectBaan:   t.User.CanSelectBaan,
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
}

func (t *UserServiceTest) TestFindOneSuccess() {
	t.User.ImageUrl = faker.URL()
	want := t.User

	c := &user.ClientMock{}
	c.On("FindOne", &proto.FindOneUserRequest{Id: t.User.Id}).Return(&proto.FindOneUserResponse{User: want}, nil)
	srv := NewService(c)

	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestFindOneNotFound() {
	want := t.NotFoundErr

	c := &user.ClientMock{}
	c.On("FindOne", &proto.FindOneUserRequest{Id: t.User.Id}).Return(nil, status.Error(codes.NotFound, "User not found"))

	srv := NewService(c)

	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *UserServiceTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	c := &user.ClientMock{}
	c.On("FindOne", &proto.FindOneUserRequest{Id: t.User.Id}).Return(nil, errors.New("Service is down"))
	srv := NewService(c)

	actual, err := srv.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *UserServiceTest) TestCreateSuccess() {
	want := t.User

	c := &user.ClientMock{}
	c.On("Create", t.UserReq).Return(&proto.CreateUserResponse{User: want}, nil)

	srv := NewService(c)

	actual, err := srv.Create(t.UserDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestCreateGrpcErr() {
	want := t.ServiceDownErr

	c := &user.ClientMock{}
	c.On("Create", t.UserReq).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.Create(t.UserDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *UserServiceTest) TestVerifySuccess() {
	c := &user.ClientMock{}
	c.On("Verify", &proto.VerifyUserRequest{StudentId: t.User.StudentID}).Return(&proto.VerifyUserResponse{Success: true}, nil)

	srv := NewService(c)

	actual, err := srv.Verify(t.User.StudentID)

	assert.Nil(t.T(), err)
	assert.True(t.T(), actual)
}

func (t *UserServiceTest) TestVerifyFailed() {
	want := t.NotFoundErr

	c := &user.ClientMock{}
	c.On("Verify", &proto.VerifyUserRequest{StudentId: t.User.StudentID}).Return(&proto.VerifyUserResponse{Success: true}, status.Error(codes.NotFound, "User not found"))

	srv := NewService(c)

	actual, err := srv.Verify(t.User.StudentID)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *UserServiceTest) TestUpdateSuccess() {
	want := t.User

	t.UserReq.Id = t.User.Id

	c := &user.ClientMock{}
	c.On("Update", t.UserReq).Return(&proto.UpdateUserResponse{User: want}, nil)

	srv := NewService(c)

	actual, err := srv.Update(t.User.Id, t.UserDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestUpdateNotFound() {
	want := t.NotFoundErr

	t.UserReq.Id = t.User.Id

	c := &user.ClientMock{}
	c.On("Update", t.UserReq).Return(nil, status.Error(codes.NotFound, "User not found"))

	srv := NewService(c)

	actual, err := srv.Update(t.User.Id, t.UserDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *UserServiceTest) TestUpdateGrpcErr() {
	want := t.ServiceDownErr

	t.UserReq.Id = t.User.Id

	c := &user.ClientMock{}
	c.On("Update", t.UserReq).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.Update(t.User.Id, t.UserDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *UserServiceTest) TestCreateOrUpdateSuccess() {
	want := t.User

	t.UserReq.Id = t.User.Id

	c := &user.ClientMock{}
	c.On("CreateOrUpdate", &proto.CreateOrUpdateUserRequest{User: t.UserReq}).Return(&proto.CreateOrUpdateUserResponse{User: t.User}, nil)

	srv := NewService(c)

	actual, err := srv.CreateOrUpdate(t.UserDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestCreateOrUpdateGrpcErr() {
	want := t.ServiceDownErr

	t.UserReq.Id = t.User.Id

	c := &user.ClientMock{}
	c.On("CreateOrUpdate", &proto.CreateOrUpdateUserRequest{User: t.UserReq}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.CreateOrUpdate(t.UserDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *UserServiceTest) TestDeleteSuccess() {
	c := &user.ClientMock{}
	c.On("Delete", &proto.DeleteUserRequest{Id: t.User.Id}).Return(&proto.DeleteUserResponse{Success: true}, nil)

	srv := NewService(c)

	actual, err := srv.Delete(t.User.Id)

	assert.Nil(t.T(), err)
	assert.True(t.T(), actual)
}

func (t *UserServiceTest) TestDeleteNotFound() {
	want := t.NotFoundErr

	c := &user.ClientMock{}
	c.On("Delete", &proto.DeleteUserRequest{Id: t.User.Id}).Return(nil, status.Error(codes.NotFound, "User not found"))

	srv := NewService(c)

	actual, err := srv.Delete(t.User.Id)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *UserServiceTest) TestDeleteGrpcErr() {
	want := t.ServiceDownErr

	c := &user.ClientMock{}
	c.On("Delete", &proto.DeleteUserRequest{Id: t.User.Id}).Return(nil, errors.New("Service is down"))

	srv := NewService(c)

	actual, err := srv.Delete(t.User.Id)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), want, err)
}
