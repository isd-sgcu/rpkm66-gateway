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

type AuthRouterTest struct {
	suite.Suite
}

func TestAuthRouter(t *testing.T) {
	suite.Run(t, new(AuthRouterTest))
}

func (t *AuthRouterTest) TestGetAuthRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "GET /auth status 200",
			route:        "/auth",
			expectedCode: http.StatusOK,
		},
		{
			description:  "GET HTTP status 404, when route is not exists",
			route:        "/not-found",
			expectedCode: http.StatusNotFound,
		},
	}

	g := mock.GuardMock{}
	conf := config.App{
		Port:        3000,
		Debug:       true,
		MaxFileSize: 1000000,
	}

	r := NewFiberRouter(&g, conf)

	r.GetAuth("/", func(ctx rctx.Context) {
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

func (t *AuthRouterTest) TestPostAuthRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "POST /auth status 201",
			route:        "/auth",
			expectedCode: http.StatusCreated,
		},
		{
			description:  "POST HTTP status 404, when route is not exists",
			route:        "/not-found",
			expectedCode: http.StatusNotFound,
		},
	}

	g := mock.GuardMock{}
	conf := config.App{
		Port:        3000,
		Debug:       true,
		MaxFileSize: 1000000,
	}

	r := NewFiberRouter(&g, conf)

	r.PostAuth("/", func(ctx rctx.Context) {
		ctx.JSON(http.StatusCreated, map[string]string{
			"message": "Hello World",
		})
	})

	for _, test := range tests {
		req := httptest.NewRequest("POST", test.route, nil)
		res, _ := r.Test(req, 1)

		assert.Equal(t.T(), test.expectedCode, res.StatusCode, test.description)
	}
}
