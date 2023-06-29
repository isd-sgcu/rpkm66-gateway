package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/pkg/rctx"
)

func (r *FiberRouter) GetAuth(path string, h func(ctx rctx.Context)) {
	r.auth.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PostAuth(path string, h func(ctx rctx.Context)) {
	r.auth.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
