package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isd-sgcu/rpkm66-gateway/src/app/handler/baan"
	"github.com/isd-sgcu/rpkm66-gateway/src/config"
	mock "github.com/isd-sgcu/rpkm66-gateway/src/mocks/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BaanRouterTest struct {
	suite.Suite
}

func TestBaanRouter(t *testing.T) {
	suite.Run(t, new(BaanRouterTest))
}

func (t *BaanRouterTest) TestPostBaanRouter() {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "POST /file status 200",
			route:        "/baan",
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

	r.GetBaan("/", func(ctx baan.IContext) {
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
