package vaccine

import (
	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	"github.com/stretchr/testify/mock"
)

type ContextMock struct {
	mock.Mock
	V      interface{}
	Status int
}

func (c *ContextMock) JSON(status int, v interface{}) {
	c.V = v
	c.Status = status
}

func (c *ContextMock) Bind(v interface{}) error {
	args := c.Called(v)

	*v.(*dto.Verify) = *args.Get(0).(*dto.Verify)

	return args.Error(1)
}

func (c *ContextMock) UserID() string {
	args := c.Called()
	return args.String(0)
}

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) Verify(hcert string, userId string) (result *dto.VaccineResponse, err *dto.ResponseErr) {
	args := s.Called(hcert, userId)

	if args.Get(0) != nil {
		result = args.Get(0).(*dto.VaccineResponse)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) Verify(req *dto.VaccineRequest, res *dto.VaccineResponse) error {
	args := c.Called(req, res)

	if args.Get(0) != nil {
		*res = *args.Get(0).(*dto.VaccineResponse)
	}

	return args.Error(1)
}
