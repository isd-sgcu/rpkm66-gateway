package vaccine

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	uMock "github.com/isd-sgcu/rpkm66-gateway/src/mocks/user"
	mock "github.com/isd-sgcu/rpkm66-gateway/src/mocks/vaccine"
	"github.com/isd-sgcu/rpkm66-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceTest struct {
	suite.Suite
	hcert          string
	User           *proto.User
	VaccineRes     *dto.VaccineResponse
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTest))
}

func (t *ServiceTest) SetupTest() {
	t.hcert = faker.Word()

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

	t.VaccineRes = &dto.VaccineResponse{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		IsPassed:  true,
		Uid:       t.User.StudentID,
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

func (t *ServiceTest) TestVerifySuccess() {
	want := t.VaccineRes

	userSrv := uMock.ServiceMock{}
	userSrv.On("Verify", t.User.StudentID, "vaccine").Return(true, nil)
	userSrv.On("FindOne", t.User.Id).Return(t.User, nil)

	req := &dto.VaccineRequest{
		HCert:     t.hcert,
		StudentId: t.User.StudentID,
	}

	c := mock.ClientMock{}
	c.On("Verify", req, &dto.VaccineResponse{}).Return(t.VaccineRes, nil)

	srv := NewService(&userSrv, &c)

	actual, err := srv.Verify(t.hcert, t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *ServiceTest) TestVerifyNotFound() {
	want := t.NotFoundErr

	userSrv := uMock.ServiceMock{}
	userSrv.On("Verify", t.User.StudentID, "vaccine").Return(false, t.NotFoundErr)
	userSrv.On("FindOne", t.User.Id).Return(t.User, nil)

	req := &dto.VaccineRequest{
		HCert:     t.hcert,
		StudentId: t.User.StudentID,
	}

	c := mock.ClientMock{}
	c.On("Verify", req, &dto.VaccineResponse{}).Return(t.VaccineRes, nil)

	srv := NewService(&userSrv, &c)

	actual, err := srv.Verify(t.hcert, t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *ServiceTest) TestVerifyInvalidQR() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Cannot verify the qr code",
		Data:       nil,
	}

	userSrv := uMock.ServiceMock{}
	userSrv.On("Verify", t.User.StudentID, "vaccine").Return(false, t.NotFoundErr)
	userSrv.On("FindOne", t.User.Id).Return(t.User, nil)

	req := &dto.VaccineRequest{
		HCert:     t.hcert,
		StudentId: t.User.StudentID,
	}

	c := mock.ClientMock{}
	c.On("Verify", req, &dto.VaccineResponse{}).Return(nil, errors.New("Invalid QR"))

	srv := NewService(&userSrv, &c)

	actual, err := srv.Verify(t.hcert, t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *ServiceTest) TestVerifyGrpcError() {
	want := t.ServiceDownErr

	userSrv := uMock.ServiceMock{}
	userSrv.On("Verify", t.User.StudentID, "vaccine").Return(false, t.ServiceDownErr)
	userSrv.On("FindOne", t.User.Id).Return(nil, t.ServiceDownErr)

	req := &dto.VaccineRequest{
		HCert:     t.hcert,
		StudentId: t.User.StudentID,
	}

	c := mock.ClientMock{}
	c.On("Verify", req, &dto.VaccineResponse{}).Return(t.VaccineRes, nil)

	srv := NewService(&userSrv, &c)

	actual, err := srv.Verify(t.hcert, t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}
