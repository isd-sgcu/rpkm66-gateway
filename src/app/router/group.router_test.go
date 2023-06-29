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

type GroupRouterTest struct {
	suite.Suite
}

func TestGroupRouter(t *testing.T) {
	suite.Run(t, new(GroupRouterTest))
}

func (t *GroupRouterTest) TestGetGroupRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "GET /Group status 200",
			route:        "/group",
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

	r.GetUser("/", func(ctx rctx.Context) {
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

func (t *GroupRouterTest) TestPostGroupRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "POST /group status 201",
			route:        "/group",
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

	r.PostGroup("/", func(ctx rctx.Context) {
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

func (t *GroupRouterTest) TestPutGroupRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "PUT /group status 200",
			route:        "/group",
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

	r.PutGroup("/", func(ctx rctx.Context) {
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

func (t *GroupRouterTest) TestDeleteGroupRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "DELETE /group status 200",
			route:        "/group",
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

	r.DeleteGroup("/", func(ctx rctx.Context) {
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
