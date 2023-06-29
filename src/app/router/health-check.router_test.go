package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isd-sgcu/rpkm66-gateway/src/config"
	mock "github.com/isd-sgcu/rpkm66-gateway/src/mocks/common"
	"github.com/isd-sgcu/rpkm66-gateway/src/pkg/rctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	conf := config.App{
		Port:        3000,
		Debug:       true,
		MaxFileSize: 1000000,
	}

	r := NewFiberRouter(&g, conf)

	r.GetHealthCheck("/", func(ctx rctx.Context) {
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
