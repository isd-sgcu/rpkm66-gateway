package vaccine

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	uMock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/user"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/vaccine"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ServiceTest struct {
	suite.Suite
	hcert          string
	studentId      string
	VaccineRes     *dto.VaccineResponse
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTest))
}

func (t *ServiceTest) SetupTest() {
	t.hcert = faker.Word()
	t.studentId = faker.Word()

	t.VaccineRes = &dto.VaccineResponse{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		IsPassed:  true,
		Uid:       t.studentId,
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
	userSrv.On("Verify", t.studentId).Return(true, nil)

	req := &dto.VaccineRequest{
		HCert:     t.hcert,
		StudentId: t.studentId,
	}

	c := mock.ClientMock{}
	c.On("Verify", req, &dto.VaccineResponse{}).Return(t.VaccineRes, nil)

	srv := NewService(&userSrv, &c)

	actual, err := srv.Verify(t.hcert, t.studentId)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *ServiceTest) TestVerifyNotFound() {
	want := t.NotFoundErr

	userSrv := uMock.ServiceMock{}
	userSrv.On("Verify", t.studentId).Return(false, t.NotFoundErr)

	req := &dto.VaccineRequest{
		HCert:     t.hcert,
		StudentId: t.studentId,
	}

	c := mock.ClientMock{}
	c.On("Verify", req, &dto.VaccineResponse{}).Return(t.VaccineRes, nil)

	srv := NewService(&userSrv, &c)

	actual, err := srv.Verify(t.hcert, t.studentId)

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
	userSrv.On("Verify", t.studentId).Return(false, t.NotFoundErr)

	req := &dto.VaccineRequest{
		HCert:     t.hcert,
		StudentId: t.studentId,
	}

	c := mock.ClientMock{}
	c.On("Verify", req, &dto.VaccineResponse{}).Return(nil, errors.New("Invalid QR"))

	srv := NewService(&userSrv, &c)

	actual, err := srv.Verify(t.hcert, t.studentId)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *ServiceTest) TestVerifyGrpcError() {
	want := t.ServiceDownErr

	userSrv := uMock.ServiceMock{}
	userSrv.On("Verify", t.studentId).Return(false, t.ServiceDownErr)

	req := &dto.VaccineRequest{
		HCert:     t.hcert,
		StudentId: t.studentId,
	}

	c := mock.ClientMock{}
	c.On("Verify", req, &dto.VaccineResponse{}).Return(t.VaccineRes, nil)

	srv := NewService(&userSrv, &c)

	actual, err := srv.Verify(t.hcert, t.studentId)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}
