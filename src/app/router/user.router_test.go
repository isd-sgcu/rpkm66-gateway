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

type UserRouterTest struct {
	suite.Suite
}

func TestUserRouter(t *testing.T) {
	suite.Run(t, new(UserRouterTest))
}

func (t *UserRouterTest) TestGetUserRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "GET /user status 200",
			route:        "/user",
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

func (t *UserRouterTest) TestPostUserRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "POST /user status 201",
			route:        "/user",
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

	r.PostUser("/", func(ctx rctx.Context) {
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

func (t *UserRouterTest) TestPutUserRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "PUT /user status 200",
			route:        "/user",
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

	r.PutUser("/", func(ctx rctx.Context) {
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

func (t *UserRouterTest) TestPatchUserRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "PATCH /user status 200",
			route:        "/user",
			expectedCode: http.StatusOK,
		},
		{
			description:  "PATCH HTTP status 404, when route is not exists",
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

	r.PatchUser("/", func(ctx rctx.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	})

	for _, test := range tests {
		req := httptest.NewRequest("PATCH", test.route, nil)
		res, _ := r.Test(req, 1)

		assert.Equal(t.T(), test.expectedCode, res.StatusCode, test.description)
	}
}

func (t *UserRouterTest) TestDeleteUserRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "DELETE /user status 200",
			route:        "/user",
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

	r.DeleteUser("/", func(ctx rctx.Context) {
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
