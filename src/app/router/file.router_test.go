package router

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/file"
	"github.com/isd-sgcu/rnkm65-gateway/src/config"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FileRouterTest struct {
	suite.Suite
}

func TestFileRouter(t *testing.T) {
	suite.Run(t, new(FileRouterTest))
}

func (t *FileRouterTest) TestPostFileRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "POST /file status 201",
			route:        "/file",
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

	r.PostFile("/", func(ctx file.IContext) {
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
