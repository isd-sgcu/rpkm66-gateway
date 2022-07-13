package router

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/user"
	"github.com/isd-sgcu/rnkm65-gateway/src/config"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type VaccineRouterTest struct {
	suite.Suite
}

func TestVaccineRouter(t *testing.T) {
	suite.Run(t, new(VaccineRouterTest))
}

func (t *VaccineRouterTest) TestGetVaccineRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "GET /vaccine status 200",
			route:        "/vaccine",
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

	r.GetMethod("/vaccine", func(ctx user.IContext) {
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

func (t *VaccineRouterTest) TestPostVaccineRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "POST /vaccine status 200",
			route:        "/vaccine",
			expectedCode: http.StatusOK,
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

	r.PostMethod("/vaccine", func(ctx user.IContext) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	})

	for _, test := range tests {
		req := httptest.NewRequest("POST", test.route, nil)
		res, _ := r.Test(req, 1)

		assert.Equal(t.T(), test.expectedCode, res.StatusCode, test.description)
	}
}

func (t *VaccineRouterTest) TestPutVaccineRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "PUT /vaccine status 200",
			route:        "/vaccine",
			expectedCode: http.StatusOK,
		},
		{
			description:  "PUT HTTP status 404, when route is not exists",
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

	r.PutMethod("/vaccine", func(ctx user.IContext) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	})

	for _, test := range tests {
		req := httptest.NewRequest("PUT", test.route, nil)
		res, _ := r.Test(req, 1)

		assert.Equal(t.T(), test.expectedCode, res.StatusCode, test.description)
	}
}

func (t *VaccineRouterTest) TestDeleteVaccineRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "DELETE /vaccine status 200",
			route:        "/vaccine",
			expectedCode: http.StatusOK,
		},
		{
			description:  "DELETE HTTP status 404, when route is not exists",
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

	r.DeleteMethod("/vaccine", func(ctx user.IContext) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	})

	for _, test := range tests {
		req := httptest.NewRequest("DELETE", test.route, nil)
		res, _ := r.Test(req, 1)

		assert.Equal(t.T(), test.expectedCode, res.StatusCode, test.description)
	}
}
