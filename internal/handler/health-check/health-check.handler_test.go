package health_check

import (
	"testing"

	"github.com/isd-sgcu/rpkm66-gateway/mocks/rctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

	c := &rctx.ContextMock{}
	h := NewHandler()

	h.HealthCheck(c)

	assert.Equal(t.T(), want, c.V)
}
