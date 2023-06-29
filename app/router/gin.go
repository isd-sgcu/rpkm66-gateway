package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/isd-sgcu/rpkm66-gateway/app/dto"
	"github.com/isd-sgcu/rpkm66-gateway/app/middleware/auth"
	"github.com/isd-sgcu/rpkm66-gateway/cfgldr"
	"github.com/isd-sgcu/rpkm66-gateway/constant/route"
	"github.com/isd-sgcu/rpkm66-gateway/pkg/rctx"
)

type GinRouter struct {
	*gin.Engine
	conf  cfgldr.App
	guard *auth.Guard
}

type Handler = func(rctx.Context) bool

func NewGinRouter(guard *auth.Guard, conf cfgldr.App) *GinRouter {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	if conf.Debug {
		// r.Use(gin.Logger())

		// add swagger
	}

	return &GinRouter{r, conf, guard}
}

func (r *GinRouter) SetHandler(key string, handler func(rctx.Context)) {
	ro, exist := route.Routes[key]

	if !exist {
		panic("Unable to find given route: " + key)
	}

	_, allPhase := ro.Phases["*"]

	if _, exist = ro.Phases[r.conf.Phase]; !exist && !allPhase {
		return
	}

	if !r.conf.Debug && ro.Debug {
		return
	}

	handlers := getMiddlewares(ro)

	ginHandler := func(ginCtx *gin.Context) {
		ctx := rctx.NewGinCtx(ginCtx)

		r.guard.Validate(ctx)

		for _, middleware := range handlers {
			if goNext := middleware(ctx); !goNext {
				return
			}
		}

		handler(ctx)
	}

	switch ro.Method {
	case route.Get:
		r.GET(ro.Path, ginHandler)
	case route.Delete:
		r.DELETE(ro.Path, ginHandler)
	case route.Patch:
		r.PATCH(ro.Path, ginHandler)
	case route.Post:
		r.POST(ro.Path, ginHandler)
	case route.Put:
		r.PUT(ro.Path, ginHandler)
	}
}

func getMiddlewares(ro route.RouteData) []Handler {
	return []Handler{getRoleMiddleware(ro.AllowPerms)}
}

func getRoleMiddleware(allowRoles map[string]struct{}) Handler {
	return func(ctx rctx.Context) bool {
		userRole := ctx.Role()

		if _, exist := allowRoles[userRole]; exist {
			return true
		} else {
			ctx.JSON(http.StatusForbidden, dto.ResponseForbiddenErr{})
			return false
		}

	}
}
