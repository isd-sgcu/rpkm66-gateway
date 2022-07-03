package router

import (
	health_check "github.com/isd-sgcu/rnkm65-gateway/src/app/handler/health-check"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HealthCheckRouterTest struct {
	suite.Suite
}

func TestHealthCheckRouter(t *testing.T) {
	suite.Run(t, new(HealthCheckRouterTest))
}

func (t *HealthCheckRouterTest) TestHealthCheckRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "Health Check with status 200",
			route:        "/",
			expectedCode: http.StatusOK,
		},
	}

	g := mock.GuardMock{}

	r := NewFiberRouter(&g, false)

	r.GetHealthCheck("/", func(ctx health_check.IContext) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	})

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)

		res, _ := r.Test(req, 1)

		assert.Equal(t.T(), test.expectedCode, res.StatusCode, test.description)
	}
}
