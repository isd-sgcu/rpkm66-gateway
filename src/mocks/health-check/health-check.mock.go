package health_check

import "github.com/stretchr/testify/mock"

type ContextMock struct {
	mock.Mock
	V interface{}
}

func (c *ContextMock) JSON(_ int, v interface{}) {
	c.V = v
}
