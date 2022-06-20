package health_check

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/mocks/health-check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HealthCheckHandlerTest struct {
	suite.Suite
}

func TestHealthCheckHandler(t *testing.T) {
	suite.Run(t, new(HealthCheckHandlerTest))
}

func (t *HealthCheckHandlerTest) TestCallHealthCheck() {
	want := map[string]interface{}{
		"Health": "OK!",
	}

	c := &health_check.ContextMock{}
	h := NewHandler()

	h.HealthCheck(c)

	assert.Equal(t.T(), want, c.V)
}
